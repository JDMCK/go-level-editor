package ui

import "github.com/hajimehoshi/ebiten/v2"

type Element interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Input interface {
	GetValue() string
	OnClick()
	OnHover()
}
