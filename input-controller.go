package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputMode int

const (
	Editing InputMode = iota
	Moving
)

type InputController struct {
	mode InputMode
}

const MoveUpKey = ebiten.KeyW
const MoveDownKey = ebiten.KeyS
const MoveLeftKey = ebiten.KeyA
const MoveRightKey = ebiten.KeyD
const FreeMoveKey = ebiten.KeySpace
const Primary = ebiten.MouseButtonLeft
const Secondary = ebiten.MouseButtonRight
const MovementSpeed = 10

func (i *InputController) Update(e *Editor) {

	kdx, kdy := handleKeyboardCameraMovement(e)
	mdx, mdy := handleMouseMovement(e)

	dx, dy := kdx+mdx, kdy+mdy

	if dx > 0 || dy > 0 || ebiten.IsKeyPressed(FreeMoveKey) {
		e.ic.mode = Moving
	} else {
		e.ic.mode = Editing
	}

	e.camera.focusX += dx
	e.camera.focusY += dy

	// mouse zoom
	_, yoff := ebiten.Wheel()
	e.camera.zoom -= yoff * 0.05
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

func CursorCamPosition(cam *Camera) (int, int) {
	mx, my := ebiten.CursorPosition()
	x := float64(mx) - cam.offsetX
	y := float64(my) - cam.offsetY

	x /= cam.zoom
	y /= cam.zoom

	x += cam.focusX
	y += cam.focusY

	return int(x), int(y)
}
