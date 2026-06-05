package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	img *ebiten.Image
}

func NewTile(img *ebiten.Image) *Tile {
	return &Tile{img: img}
}

func NewEmptyTile(tileWidth, tileHeight int) *Tile {
	return &Tile{img: ebiten.NewImage(tileWidth, tileHeight)}
}

func (t *Tile) Draw(x, y int, screen *ebiten.Image, cam *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.GeoM.Concat(cam.GeoM)
	screen.DrawImage(t.img, op)
}

// width in tiles (not pixels)
func TilePositionFromIndex(index int, width int, tileWidth, tileHeight int) (int, int) {
	return (index % width) * tileWidth, (index / width) * tileHeight
}

func TileIndexFromPosition(x, y int, width, height int, tileWidth, tileHeight int) int {
	if x < 0 || y < 0 {
		return -1
	}
	sX := x / tileWidth
	sY := y / tileHeight

	if sX >= width || sY >= height {
		return -1
	}

	return (sY * width) + sX
}
