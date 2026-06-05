package main

import (
	"image/color"
	gui "level-editor/gui"

	"github.com/hajimehoshi/ebiten/v2"
	eb "github.com/hajimehoshi/ebiten/v2"
)

const TileSize = 16

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

	ic InputController

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

	e.camera = NewCamera()
	e.camera.CenterScreenOffset(ScreenWidth, ScreenHeight)

	layers := make([]*Container, 0)
	layers = append(layers, NewEmptyContainer(0, 0, TileSize, TileSize, 50, 25))
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

	// update cursor (draw / erase)
	curTile, x, y := e.canvas[e.currLayer].TileFromCursor(e.camera)
	e.cursor.SelectTile(x, y, curTile)
	if e.cursor.tile != nil && ebiten.IsMouseButtonPressed(Primary) && e.ic.mode == Editing {
		e.cursor.GetTile().img.Fill(color.RGBA{255, 255, 0, 255})
	}
	if e.cursor.tile != nil && ebiten.IsMouseButtonPressed(Secondary) && e.ic.mode == Editing {
		e.cursor.GetTile().img.Clear()
	}

	e.palette.Update()

	return nil
}

func (e *Editor) Draw(screen *eb.Image) {
	for _, l := range e.canvas {
		l.Draw(screen, e.camera.DrawOptions())
	}

	for _, el := range e.GUI {
		el.Draw(screen)
	}
	e.cursor.Draw(screen, e.camera.DrawOptions())

	e.palette.Draw(screen)
}

func (e *Editor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
