package main

import (
	"fmt"
	"image/color"
	gui "level-editor/gui"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/sqweek/dialog"
)

const TileSize = 16
const CanvasWidth = 25
const CanvasHeight = 25

var ScreenWidth = 1920
var ScreenHeight = 1080

const DefaultLayerCount = 3
const MinLayerCount = 1
const MaxLayerCount = 10

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

	showGUI bool
}

type Drawable interface {
	Draw(dst *eb.Image)
}

func NewEditor() *Editor {
	e := Editor{}

	e.camera = NewCamera(3)
	e.camera.CenterScreenOffset(ScreenWidth, ScreenHeight)
	e.camera.focusX += CanvasWidth * TileSize / 2
	e.camera.focusY += CanvasHeight * TileSize / 2

	layers := buildCanvas()
	e.canvas = layers
	e.layerVisibility = make([]bool, MaxLayerCount)

	cursor := NewCursor(TileSize, TileSize)
	e.cursor = cursor

	// p := NewPaletteFromTileMap(10, 10, "assets/dungeon.png", 16, 16, 6, 18)
	p := NewPaletteFromTileMap(10, 10, "assets/basictiles.png", 16, 16, 8, 120)
	e.palette = p

	e.GUI = buildGUI(&e)
	e.showGUI = true

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

	if e.showGUI {
		for _, el := range e.GUI {
			el.Draw(screen)
		}
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
	if eb.IsMouseButtonPressed(Primary) {
		tile := e.palette.SelectedTile()
		if tile != nil && curTile != nil {
			curTile.img = tile.img
		}
	}
	if eb.IsMouseButtonPressed(Secondary) {
		e.cursor.Tile.Reset()
	}
}
func handleBlockEditMode(e *Editor) {
	// TODO
}

func handleMovingMode(e *Editor) {
	kdx, kdy := handleKeyboardCameraMovement(e)
	mdx, mdy := handleMouseMovement(e)

	dx, dy := kdx+mdx, kdy+mdy

	e.camera.focusX += dx
	e.camera.focusY += dy
}

func handleZoom(e *Editor) {
	_, yoff := ebiten.Wheel()
	e.camera.zoom -= yoff * ZoomSpeed
}

func buildCanvas() []*Container {
	layers := make([]*Container, 0, MaxLayerCount)
	for range DefaultLayerCount {
		layers = append(layers, NewEmptyContainer(0, 0, TileSize, TileSize, CanvasWidth, CanvasHeight))
	}
	return layers
}

func buildGUI(e *Editor) []gui.Element {
	rootX, rootY := ScreenWidth-250, 50
	gap := 60
	guiEls := make([]gui.Element, 0)
	// guiEls = append(guiEls, gui.NewBasicButton("Save", 50, 50, gui.Large, gui.Primary, func() { fmt.Println("Saved") }))
	// guiEls = append(guiEls, gui.NewBasicButton("Cancel", 50, 100, gui.Small, gui.Secondary, func() { fmt.Println("Cancel") }))
	// guiEls = append(guiEls, gui.NewBasicButton("Delete", 50, 150, gui.Medium, gui.Danger, func() { fmt.Println("Delete") }))
	// guiEls = append(guiEls, gui.NewNumberPicker(1, 5, 0, 50, 200))
	// guiEls = append(guiEls, gui.NewCheckbox(50, 250))
	clearCanvas := func() {
		e.canvas[e.currLayer].Clear()
	}
	layerCountChange := func(val int) {
		if val > MaxLayerCount || val < MinLayerCount || val == len(e.canvas) {
			return
		}
		if val > len(e.canvas) {
			// diff between val and canvas len can only ever be max 1
			e.canvas = append(e.canvas, NewEmptyContainer(0, 0, TileSize, TileSize, CanvasWidth, CanvasHeight))
			e.layerVisibility[len(e.canvas)-1] = true
		}
		if len(e.canvas) > val {
			e.canvas = e.canvas[:val]
		}
		e.currLayer = val - 1
		e.GUI = buildGUI(e) // rebuild gui to have more layer buttons
	}
	saveLevel := newSaveAction()
	toggleLayer := func(visible bool, i int) {
		e.layerVisibility[i] = visible
	}
	layerChange := func(i int) {
		e.currLayer = i
		e.GUI = buildGUI(e) // rebuild gui to update current layer label (terribly inefficient, I know)
	}

	guiEls = append(guiEls, gui.NewBasicButton("Save", rootX, rootY, gui.Large, gui.Primary, saveLevel))
	guiEls = append(guiEls, gui.NewBasicButton("Clear Layer", rootX, rootY+gap, gui.Large, gui.Danger, clearCanvas))
	guiEls = append(guiEls, gui.NewNumberPicker("Layer Count", len(e.canvas), MinLayerCount, MaxLayerCount, rootX, rootY+2*gap, layerCountChange))
	guiEls = append(guiEls, gui.NewText(fmt.Sprintf("Current Layer: %d", e.currLayer+1), rootX, rootY+3*gap, 20, color.White))
	for i := range len(e.canvas) {
		guiEls = append(guiEls, gui.NewCheckbox(rootX, rootY+(i+4)*gap, true, func(visible bool) { toggleLayer(visible, i) }))
		guiEls = append(guiEls, gui.NewBasicButton(fmt.Sprintf("Edit Layer %d", i+1), rootX+gap, rootY+(i+4)*gap, gui.Medium, gui.Secondary, func() { layerChange(i) }))
	}

	return guiEls
}

func newSaveAction() func() {
	return func() {
		file, err := dialog.
			File().
			Title("Save Level").
			SetStartFile("level00.config.map").
			Filter("map", "txt").
			Save()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Selected file:", file)
	}
}
