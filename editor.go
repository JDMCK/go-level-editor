package main

import (
	gui "level-editor/gui"

	"github.com/hajimehoshi/ebiten/v2"
	eb "github.com/hajimehoshi/ebiten/v2"
)

const TileSize = 16
const CanvasWidth = 25
const CanvasHeight = 25

var ScreenWidth = 1920
var ScreenHeight = 1080

type Editor struct {
	camera *Camera
	GUI    []gui.Element

	prevCursorX int
	prevCursorY int

	palette   *Palette
	canvas    []*Container
	currLayer int

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

	layers := make([]*Container, 0)
	layers = append(layers, NewEmptyContainer(0, 0, TileSize, TileSize, CanvasWidth, CanvasHeight))
	e.canvas = layers

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
	for _, l := range e.canvas {
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

func buildGUI(e *Editor) []gui.Element {
	rootX, rootY := ScreenWidth-250, 50
	guiEls := make([]gui.Element, 0)
	// guiEls = append(guiEls, gui.NewBasicButton("Save", 50, 50, gui.Large, gui.Primary, func() { fmt.Println("Saved") }))
	// guiEls = append(guiEls, gui.NewBasicButton("Cancel", 50, 100, gui.Small, gui.Secondary, func() { fmt.Println("Cancel") }))
	// guiEls = append(guiEls, gui.NewBasicButton("Delete", 50, 150, gui.Medium, gui.Danger, func() { fmt.Println("Delete") }))
	// guiEls = append(guiEls, gui.NewNumberPicker(1, 5, 0, 50, 200))
	// guiEls = append(guiEls, gui.NewCheckbox(50, 250))
	clearCanvas := func() {
		e.canvas[e.currLayer].Clear()
	}

	layerChange := func(val int) {
		if val < len(e.canvas) && val > 0 {
			e.currLayer = val
		}
	}

	guiEls = append(guiEls, gui.NewBasicButton("Clear", rootX, rootY, gui.Medium, gui.Danger, clearCanvas))
	guiEls = append(guiEls, gui.NewNumberPicker("Layer", 0, 0, 3, rootX, rootY+50, layerChange))

	return guiEls
}
