package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputMode int

const (
	Editing InputMode = iota
	BlockEdit
	Moving
)

type InputController struct {
	Mode InputMode
}

const MoveUpKey = ebiten.KeyW
const MoveDownKey = ebiten.KeyS
const MoveLeftKey = ebiten.KeyA
const MoveRightKey = ebiten.KeyD
const FreeMoveKey = ebiten.KeySpace
const Primary = ebiten.MouseButtonLeft
const Secondary = ebiten.MouseButtonRight
const MultiSelect = ebiten.KeyControl // TODO
const MovementSpeed = 10
const ZoomSpeed = 0.1

func (i *InputController) Update(e *Editor) {
	switch {
	case isMovementMode():
		e.ic.Mode = Moving
	case ebiten.IsKeyPressed(MultiSelect):
		e.ic.Mode = BlockEdit
	default:
		e.ic.Mode = Editing
	}
}

func isMovementMode() bool {
	return ebiten.IsKeyPressed(FreeMoveKey) ||
		ebiten.IsKeyPressed(MoveUpKey) ||
		ebiten.IsKeyPressed(MoveDownKey) ||
		ebiten.IsKeyPressed(MoveLeftKey) ||
		ebiten.IsKeyPressed(MoveRightKey)
}

func handleMouseMovement(e *Editor) (float64, float64) { // mouse drag
	x, y := ebiten.CursorPosition()
	var dx, dy float64
	if ebiten.IsKeyPressed(FreeMoveKey) && ebiten.IsMouseButtonPressed(Primary) {
		dx = float64(e.prevCursorX - x)
		dx /= e.camera.zoom
		dy = float64(e.prevCursorY - y)
		dy /= e.camera.zoom
	}
	e.prevCursorX = x
	e.prevCursorY = y

	return dx, dy
}

func handleKeyboardCameraMovement(e *Editor) (float64, float64) {
	var dx float64 = 0
	var dy float64 = 0
	if ebiten.IsKeyPressed(MoveUpKey) {
		dy += -MovementSpeed / e.camera.zoom
	}
	if ebiten.IsKeyPressed(MoveDownKey) {
		dy += MovementSpeed / e.camera.zoom
	}
	if ebiten.IsKeyPressed(MoveLeftKey) {
		dx += -MovementSpeed / e.camera.zoom
	}
	if ebiten.IsKeyPressed(MoveRightKey) {
		dx += MovementSpeed / e.camera.zoom
	}
	return dx, dy
}

func CursorPosition(op *ebiten.DrawImageOptions) (int, int) {
	mx, my := ebiten.CursorPosition()
	op.GeoM.Invert()
	x, y := op.GeoM.Apply(float64(mx), float64(my))

	return int(x), int(y)
}
