package main

import (
	"fmt"
	gui "level-editor/gui"

	"github.com/hajimehoshi/ebiten/v2"
	eb "github.com/hajimehoshi/ebiten/v2"
)

type Editor struct {
	camera *Camera
	GUI    []gui.Element

	prevCursorX int
	prevCursorY int
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
	return &e
}

func (e *Editor) Update() error {
	for _, el := range e.GUI {
		el.Update()
	}

	if ebiten.IsKeyPressed(FreeMoveKey) && ebiten.IsMouseButtonPressed(Primary) {
		x, y := ebiten.CursorPosition()
		dx := x - e.prevCursorX
		dy := y - e.prevCursorY
		e.camera.focusX += float64(dx)
		e.camera.focusY += float64(dy)
		return nil
	}
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
	for _, el := range e.GUI {
		el.Draw(screen)
	}
}

func (e *Editor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800
}
