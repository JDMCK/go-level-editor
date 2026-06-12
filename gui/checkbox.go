package gui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Checkbox struct {
	checked bool
	btn     *Button
	onCheck func(bool)
}

func NewCheckbox(x, y int, checked bool, onCheck func(bool)) *Checkbox {
	cb := Checkbox{
		checked: checked,
		onCheck: onCheck,
	}
	btn := NewButton("", 32, 32, x, y, 0, color.Gray{150}, func() {
		cb.checked = !cb.checked
	})
	cb.btn = btn
	return &cb
}

func (c *Checkbox) Update() {
	if c.checked {
		c.btn.SetColor(color.RGBA{0, 255, 0, 255})
		c.onCheck(true)
	} else {
		c.btn.SetColor(color.Gray{150})
		c.onCheck(false)
	}
	c.btn.Update()
}

func (c *Checkbox) Draw(dst *ebiten.Image) {
	c.btn.Draw(dst)
}
