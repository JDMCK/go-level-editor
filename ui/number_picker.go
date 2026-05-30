package ui

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type NumberPicker struct {
	x, y    int
	value   int
	buttons [4]*Button
}

const defaultNPButtonSize = 18

var defaultNPButtonColor = color.Gray{100}

func newNumberPickerButton(jump int, x, y int, onClick func()) *Button {
	label := strconv.Itoa(jump)
	return NewButton(label, defaultNPButtonSize, defaultNPButtonSize, x, y, 16, defaultNPButtonColor, onClick)
}

func NewNumberPicker(smallJump, bigJump int, initialValue int, x, y int) *NumberPicker {
	np := NumberPicker{
		value: initialValue,
		x:     x,
		y:     y,
	}
	btns := [4]*Button{
		newNumberPickerButton(-bigJump, x, y, func() { np.value -= bigJump }),
		newNumberPickerButton(-smallJump, x+defaultNPButtonSize, y, func() { np.value -= smallJump }),
		newNumberPickerButton(smallJump, x+(3*defaultNPButtonSize), y, func() { np.value += smallJump }),
		newNumberPickerButton(bigJump, x+(4*defaultNPButtonSize), y, func() { np.value += bigJump }),
	}
	np.buttons = btns
	return &np
}

func (n *NumberPicker) Update() {

}

func (n *NumberPicker) Draw(screen *ebiten.Image) {
	for _, btn := range n.buttons {
		btn.Draw(screen)
	}
	valueBox := ebiten.NewImage(defaultNPButtonSize, defaultNPButtonSize)
	valueBox.Fill(defaultNPButtonColor)
	text := NewText(strconv.Itoa(n.value), n.x+(2*defaultNPButtonSize), n.y, 16, defaultNPButtonColor)
	screen.DrawImage(valueBox, nil)
	text.Draw(screen)
}
