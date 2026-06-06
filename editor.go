package main

import (
	"fmt"
	gui "level-editor/gui"

	"github.com/hajimehoshi/ebiten/v2"
	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

	multiStartTile *Tile

	palette   *Palette
	canvas    []*Container
	currLayer int

	ic     InputController
	cursor *Cursor
}

type Drawable interface {
	Draw(dst *ebiten.Image)
}

func NewEditor() *Editor {
	e := Editor{}

	// guiEls := make([]gui.Element, 0)
	// guiEls = append(guiEls, gui.NewBasicButton("Save", 50, 50, gui.Large, gui.Primary, func() { fmt.Println("Saved") }))
	// guiEls = append(guiEls, gui.NewBasicButton("Cancel", 50, 100, gui.Small, gui.Secondary, func() { fmt.Println("Cancel") }))
	// guiEls = append(guiEls, gui.NewBasicButton("Delete", 50, 150, gui.Medium, gui.Danger, func() { fmt.Println("Delete") }))
	// guiEls = append(guiEls, gui.NewNumberPicker(1, 5, 0, 50, 200))
	// guiEls = append(guiEls, gui.NewCheckbox(50, 250))
	// e.GUI = guiEls

	e.camera = NewCamera(3)
	e.camera.CenterScreenOffset(ScreenWidth, ScreenHeight)
	e.camera.focusX += CanvasWidth * TileSize / 2
	e.camera.focusY += CanvasHeight * TileSize / 2

	layers := make([]*Container, 0)
	layers = append(layers, NewEmptyContainer(0, 0, TileSize, TileSize, CanvasWidth, CanvasHeight))
	e.canvas = layers

	cursor := NewCursor(TileSize, TileSize)
	e.cursor = cursor

	p := NewPaletteFromTileMap(10, 10, "assets/dungeon.png", 16, 16, 6, 18)
	e.palette = p

	return &e
}

func (e *Editor) Update() error {
	for _, el := range e.GUI {
		el.Update()
	}
	e.ic.Update(e)
	e.palette.Update()

	camOp := e.camera.DrawOptions()

	curTile := e.canvas[e.currLayer].TileFromCursor(camOp)
	e.cursor.SelectTile(curTile)

	// draw / erase
	if e.cursor.tile == nil {
		return nil
	}
	if ebiten.IsMouseButtonPressed(Primary) && e.ic.mode == Editing {
		tile := e.palette.SelectedTile()
		if tile != nil {
			e.cursor.tile.img = tile.img
		}
	}
	if ebiten.IsMouseButtonPressed(Secondary) && e.ic.mode == Editing {
		e.cursor.tile.Reset()
	}
	// block draw / erase
	if inpututil.IsMouseButtonJustPressed(Primary) && e.ic.mode == BlockEdit {
		e.multiStartTile = curTile
	}
	if e.ic.mode != BlockEdit {
		e.multiStartTile = nil
	}
	if e.ic.mode == BlockEdit && e.multiStartTile == nil {
		e.cursor.SelectTile(curTile)
	}
	if e.multiStartTile != nil && e.ic.mode == BlockEdit {
		midX := curTile.x
		midY := curTile.y

		topRightX := min(e.multiStartTile.x, midX)
		topRightY := max(e.multiStartTile.y, midY)

		botRightX := max(e.multiStartTile.x, midX)
		botLeftY := max(e.multiStartTile.y, midY)

		topRightTile := e.canvas[e.currLayer].TileFromPosition(topRightX, topRightY, camOp)
		fmt.Println(topRightTile)
		e.cursor.MultiSelect(topRightTile, botRightX, botLeftY)
	}

	return nil
}

func (e *Editor) Draw(screen *eb.Image) {
	for _, l := range e.canvas {
		l.Draw(screen, e.camera.DrawOptions())
	}
	e.cursor.Draw(screen, e.camera.DrawOptions())

	for _, el := range e.GUI {
		el.Draw(screen)
	}

	e.palette.Draw(screen)
}

func (e *Editor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
