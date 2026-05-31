package main

import "github.com/hajimehoshi/ebiten/v2"

type Layer struct {
	tiles  []Tile
	width  int // in tiles
	height int // in tiles
}

func (l *Layer) Draw(screen *ebiten.Image, cam *ebiten.DrawImageOptions) {
	for i, t := range l.tiles {
		x, y := TilePosition(i, l.width)
		t.Draw(x, y, screen, cam)
	}
}

func TilePosition(index int, width int) (int, int) {
	return index % width, index / width
}
