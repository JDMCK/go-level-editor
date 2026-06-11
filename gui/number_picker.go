package gui

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

type NumberPicker struct {
	x, y     int
	value    int
	buttons  [2]*Button
	valueImg *ebiten.Image
	rect     Rectangle
	text     Text
	minValue int
	maxValue int
	label    Text
}

const defaultNPButtonSize = 18

var defaultNPButtonColor = color.Gray{100}
var defaultNPLabelColor = color.White
var defaultNPLabelSize = 16.0

func newNumberPickerButton(label string, x, y int, onClick func()) *Button {
	return NewButton(label, defaultNPButtonSize, defaultNPButtonSize, x, y, 16, defaultNPButtonColor, onClick)
}

func NewNumberPicker(label string, initialValue, minValue, maxValue int, x, y int, onChange func(int)) *NumberPicker {
	labelTxt := NewText(label, x, y, defaultNPLabelSize, defaultNPButtonColor)
	np := NumberPicker{
		value: initialValue,
		x:     x,
		y:     y,
		label: labelTxt,
	}
	btns := [2]*Button{
		newNumberPickerButton("-", x, y+labelTxt.rect.height, func() {
			if np.value-1 >= minValue {
				np.value -= 1
			}
			onChange(np.value)
		}),
		newNumberPickerButton("+", x+(3*defaultNPButtonSize), y+labelTxt.rect.height, func() {
			if np.value+1 <= maxValue {
				np.value += 1
			}
			onChange(np.value)
		}),
	}
	np.buttons = btns
	np.rect = Rectangle{x, y + labelTxt.rect.height, defaultNPButtonSize * 4, defaultNPButtonSize}
	text := NewText(strconv.Itoa(np.value), np.x, np.y, defaultNPLabelSize, defaultNPButtonColor)
	text.CenterInRectangle(np.rect)
	np.text = text
	return &np
}

func (n *NumberPicker) Update() {
	for _, el := range n.buttons {
		el.Update()
	}
	n.text.SetValue(strconv.Itoa(n.value))
	n.text.CenterInRectangle(n.rect)
}

func (n *NumberPicker) Draw(screen *ebiten.Image) {
	n.label.Draw(screen)
	for _, btn := range n.buttons {
		btn.Draw(screen)
	}
	n.text.Draw(screen)
}
