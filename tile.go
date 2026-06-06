package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	x, y int
	img  *ebiten.Image
}

func NewTile(x, y int, img *ebiten.Image) *Tile {
	return &Tile{x, y, img}
}

func NewEmptyTile(x, y int, tileWidth, tileHeight int) *Tile {
	return &Tile{x: x, y: y, img: ebiten.NewImage(tileWidth, tileHeight)}
}

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

func (t *Tile) Reset() {
	p := t.img.Bounds().Size()
	t.img = ebiten.NewImage(p.X, p.Y)
}

func (t *Tile) Draw(screen *ebiten.Image, cam *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.x), float64(t.y))
	op.GeoM.Concat(cam.GeoM)
	screen.DrawImage(t.img, op)
}
