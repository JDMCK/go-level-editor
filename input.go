package main

import "github.com/hajimehoshi/ebiten/v2"

func CursorPosition(cam *Camera) (int, int) {
	mx, my := ebiten.CursorPosition()
	return mx - int(cam.focusX) + int(cam.offsetX), my - int(cam.focusY) + int(cam.offsetY)
}
