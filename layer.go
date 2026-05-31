package main

import "github.com/hajimehoshi/ebiten/v2"

type Layer struct {
	tiles []Tile
	width int // in tiles
}

func (l *Layer) Draw(screen *ebiten.Image, cam *ebiten.DrawImageOptions) {
	for i, t := range l.tiles {
		x, y := TilePosition(i, l.width)
		t.Draw(x, y, screen, cam)
	}
}
