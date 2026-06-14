package main

import (
	"fmt"
	"image/color"
	gui "level-editor/gui"

	"github.com/hajimehoshi/ebiten/v2"
	eb "github.com/hajimehoshi/ebiten/v2"
)

var ScreenWidth = 1920
var ScreenHeight = 1080

const DefaultLayerCount = 3
const MinLayerCount = 1
const MaxLayerCount = 10

var CanvasWidth int = 32
var CanvasHeight int = 16

type Editor struct {
	camera *Camera
	GUI    []gui.Element

	prevCursorX int
	prevCursorY int

	palette         *Palette
	canvas          []*Container
	currLayer       int
	layerVisibility []bool

	ic     InputController
	cursor *Cursor
}

type Drawable interface {
	Draw(dst *eb.Image)
}

func NewEditor() *Editor {
	e := Editor{}

	e.camera = NewCamera(3)
	e.camera.CenterScreenOffset(ScreenWidth, ScreenHeight)
	e.camera.focusX += float64(CanvasWidth * TileWidth / 2)
	e.camera.focusY += float64(CanvasHeight * TileHeight / 2)

	cursor := NewCursor(TileWidth, TileHeight)
	e.cursor = cursor

	p := NewPaletteFromTileMap(10, 10, *AtlasPath)
	e.palette = p

	layers := buildCanvas(p)
	e.canvas = layers
	e.layerVisibility = make([]bool, MaxLayerCount)

	e.GUI = buildGUI(&e)

	return &e
}

func (e *Editor) Update() error {
	for _, el := range e.GUI {
		el.Update()
	}
	e.ic.Update(e)
	e.palette.Update()

	// draw / erase
	switch e.ic.Mode {
	case Editing:
		handleEditMode(e)
	case BlockEdit:
		handleBlockEditMode(e)
	case Moving:
		handleMovingMode(e)
	}
	handleZoom(e)

	return nil
}

func (e *Editor) Draw(screen *eb.Image) {
	for i, l := range e.canvas {
		if e.layerVisibility[i] == false {
			continue
		}
		l.Draw(screen, e.camera.DrawOptions())
	}
	e.cursor.Draw(screen, e.camera.DrawOptions())

	e.palette.Draw(screen)

	for _, el := range e.GUI {
		el.Draw(screen)
	}
}

func (e *Editor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if outsideWidth != ScreenWidth || outsideHeight != ScreenHeight {
		ScreenWidth = outsideWidth
		ScreenHeight = outsideHeight
		e.GUI = buildGUI(e)
	}
	return outsideWidth, outsideHeight
}

func handleEditMode(e *Editor) {
	camOp := e.camera.DrawOptions()
	curTile := e.canvas[e.currLayer].TileFromCursor(camOp)
	e.cursor.SelectTile(curTile)

	if e.cursor.Tile == nil {
		return
	}

	handleEdits(e, curTile)
}

func handleBlockEditMode(e *Editor) {
	// TODO
}

func handleMovingMode(e *Editor) {
	dx, dy := handleMouseMovement(e)
	e.camera.focusX += dx
	e.camera.focusY += dy
}

func handleZoom(e *Editor) {
	_, yoff := ebiten.Wheel()
	e.camera.zoom -= yoff * ZoomSpeed
}

func buildCanvas(p *Palette) []*Container {
	importedLayersLen := len(ImportedLayerIndices)
	if importedLayersLen != 0 {
		layers := make([]*Container, 0, importedLayersLen)
		for i := range importedLayersLen {
			layers = append(layers, NewContainerFromAtlasIndices(0, 0, CanvasWidth, CanvasHeight, i, p))
		}
		return layers
	}
	layers := make([]*Container, 0, MaxLayerCount)
	for range DefaultLayerCount {
		layers = append(layers, NewEmptyContainer(0, 0, CanvasWidth, CanvasHeight))
	}
	return layers
}

func buildGUI(e *Editor) []gui.Element {
	rootX, rootY := ScreenWidth-250, 50
	gap := 60
	guiEls := make([]gui.Element, 0)

	// callback methods
	clearCanvas := func() {
		e.canvas[e.currLayer].Clear()
	}
	layerCountChange := func(val int) {
		if val > MaxLayerCount || val < MinLayerCount || val == len(e.canvas) {
			return
		}
		if val > len(e.canvas) {
			// diff between val and canvas len can only ever be max 1
			e.canvas = append(e.canvas, NewEmptyContainer(0, 0, CanvasWidth, CanvasHeight))
			e.layerVisibility[len(e.canvas)-1] = true
		}
		if len(e.canvas) > val {
			e.canvas = e.canvas[:val]
		}
		e.currLayer = val - 1
		e.GUI = buildGUI(e) // rebuild gui to have more layer buttons
	}
	saveLevel := NewSaveAction(e)
	toggleLayer := func(visible bool, i int) {
		e.layerVisibility[i] = visible
	}
	layerChange := func(i int) {
		e.currLayer = i
		e.GUI = buildGUI(e) // rebuild gui to update current layer label (terribly inefficient, I know)
	}
	canvasWidthChange := func(val int) {
		for _, l := range e.canvas {
			l.SetWidth(val)
		}
	}
	canvasHeightChange := func(val int) {
		for _, l := range e.canvas {
			l.SetHeight(val)
		}
	}

	guiEls = append(guiEls, gui.NewBasicButton("Save", rootX, rootY, gui.Large, gui.Primary, saveLevel))
	guiEls = append(guiEls, gui.NewBasicButton("Clear Layer", rootX, rootY+gap, gui.Large, gui.Danger, clearCanvas))
	guiEls = append(guiEls, gui.NewNumberPicker("Layer Count", len(e.canvas), MinLayerCount, MaxLayerCount, rootX, rootY+2*gap, layerCountChange))
	guiEls = append(guiEls, gui.NewText(fmt.Sprintf("Current Layer: %d", e.currLayer+1), rootX, rootY+3*gap, 20, color.White))
	for i := range len(e.canvas) {
		guiEls = append(guiEls, gui.NewCheckbox(rootX, rootY+(i+4)*gap, true, func(visible bool) { toggleLayer(visible, i) }))
		guiEls = append(guiEls, gui.NewBasicButton(fmt.Sprintf("Edit Layer %d", i+1), rootX+gap, rootY+(i+4)*gap, gui.Medium, gui.Secondary, func() { layerChange(i) }))
	}
	guiEls = append(guiEls, gui.NewNumberPicker("Canvas Width", CanvasWidth, 1, 300, rootX+2*gap, ScreenHeight-3*gap, canvasWidthChange))
	guiEls = append(guiEls, gui.NewNumberPicker("Canvas Height", CanvasHeight, 1, 300, rootX+2*gap, ScreenHeight-2*gap, canvasHeightChange))

	return guiEls
}
