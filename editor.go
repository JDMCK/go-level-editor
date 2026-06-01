package main

import (
	"fmt"
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
	layers    []*Container
	currLayer int
}

type Drawable interface {
	Draw(dst *ebiten.Image)
}

func NewEditor() *Editor {
	e := Editor{}

	guiEls := make([]gui.Element, 0)
	guiEls = append(guiEls, gui.NewBasicButton("Save", 50, 50, gui.Large, gui.Primary, func() { fmt.Println("Saved") }))
	guiEls = append(guiEls, gui.NewBasicButton("Cancel", 50, 100, gui.Small, gui.Secondary, func() { fmt.Println("Cancel") }))
	guiEls = append(guiEls, gui.NewBasicButton("Delete", 50, 150, gui.Medium, gui.Danger, func() { fmt.Println("Delete") }))
	guiEls = append(guiEls, gui.NewNumberPicker(1, 5, 0, 50, 200))
	guiEls = append(guiEls, gui.NewCheckbox(50, 250))
	e.GUI = guiEls

	e.camera = NewCamera()
	e.camera.CenterScreenOffset(ScreenWidth, ScreenHeight)

	layers := make([]*Container, 0)
	layers = append(layers, NewEmptyContainer(0, 0, TileSize, TileSize, 25, 25))
	e.layers = layers
	return &e
}

func (e *Editor) Update() error {
	for _, el := range e.GUI {
		el.Update()
	}

	// cursor drag
	x, y := ebiten.CursorPosition()
	if ebiten.IsKeyPressed(FreeMoveKey) && ebiten.IsMouseButtonPressed(Primary) {
		dx := float64(e.prevCursorX - x)
		dx *= 1 / e.camera.zoom
		dy := float64(e.prevCursorY - y)
		dy *= 1 / e.camera.zoom
		e.camera.focusX += dx
		e.camera.focusY += dy
	}
	e.prevCursorX = x
	e.prevCursorY = y

	// cursor zoom
	_, yoff := ebiten.Wheel()
	e.camera.zoom -= yoff * 0.01

	fmt.Println(e.camera)

	handleKeyboardCameraMovement(e)
	return nil
}

func handleKeyboardCameraMovement(e *Editor) {
	var velX float64 = 0
	var velY float64 = 0

	if ebiten.IsKeyPressed(MoveUpKey) {
		velY += -MovementSpeed
	}
	if ebiten.IsKeyPressed(MoveDownKey) {
		velY += MovementSpeed
	}
	if ebiten.IsKeyPressed(MoveLeftKey) {
		velX += -MovementSpeed
	}
	if ebiten.IsKeyPressed(MoveRightKey) {
		velX += MovementSpeed
	}

	e.camera.focusX += velX
	e.camera.focusY += velY
}

func (e *Editor) Draw(screen *eb.Image) {
	for _, l := range e.layers {
		l.Draw(screen, e.camera.DrawOptions())
	}

	for _, el := range e.GUI {
		el.Draw(screen)
	}
}

func (e *Editor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
