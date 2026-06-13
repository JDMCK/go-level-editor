package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

const MovePaletteCursorUpKey = ebiten.KeyW
const MovePaletteCursorDownKey = ebiten.KeyS
const MovePaletteCursorLeftKey = ebiten.KeyA
const MovePaletteCursorRightKey = ebiten.KeyD
const FreeMoveKey = ebiten.KeySpace
const Primary = ebiten.MouseButtonLeft
const Secondary = ebiten.MouseButtonRight
const Tertiary = ebiten.MouseButtonMiddle
const MultiSelect = ebiten.KeyControl // TODO
const MovementSpeed = 10
const ZoomSpeed = 0.1

func (i *InputController) Update(e *Editor) {
	switch {
	case ebiten.IsKeyPressed(FreeMoveKey):
		e.ic.Mode = Moving
	case ebiten.IsKeyPressed(MultiSelect):
		e.ic.Mode = BlockEdit
	default:
		e.ic.Mode = Editing
	}
	handleKeyboardCursorMovement(e.palette)
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

// Keyboard movement for palette
func handleKeyboardCursorMovement(p *Palette) {
	c := p.selectionCursor
	pWidth := p.container.width
	pHeight := p.container.height
	if c.Tile == nil {
		return
	}

	cIndex := TileIndexFromPosition(c.Tile.x, c.Tile.y, pWidth, pHeight)
	i := -1

	switch {
	case inpututil.IsKeyJustPressed(MovePaletteCursorUpKey):
		if cIndex < pWidth {
			i = cIndex + pWidth*(pHeight-1)
			fmt.Println(i, cIndex)
		} else {
			i = cIndex - pWidth
		}
	case inpututil.IsKeyJustPressed(MovePaletteCursorDownKey):
		if cIndex/pWidth == pHeight-1 {
			i = cIndex - pWidth*(pHeight-1)
		} else {
			i = cIndex + pWidth
		}
	case inpututil.IsKeyJustPressed(MovePaletteCursorLeftKey):
		if cIndex%pWidth == 0 {
			i = cIndex + pWidth - 1
		} else {
			i = cIndex - 1
		}
	case inpututil.IsKeyJustPressed(MovePaletteCursorRightKey):
		if cIndex%pWidth == pWidth-1 {
			i = cIndex - pWidth + 1
		} else {
			i = cIndex + 1
		}
	}
	if i == -1 {
		return
	}
	c.SelectTile(p.container.tiles[i])
}

func CursorPosition(op *ebiten.DrawImageOptions) (int, int) {
	mx, my := ebiten.CursorPosition()
	op.GeoM.Invert()
	x, y := op.GeoM.Apply(float64(mx), float64(my))

	return int(x), int(y)
}
