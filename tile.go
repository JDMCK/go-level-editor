package main

import "github.com/hajimehoshi/ebiten/v2"

type Tile struct {
	img *ebiten.Image
}

func NewTile(img *ebiten.Image) *Tile {
	return &Tile{img: img}
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

func TileIndexFromPosition(x, y int, width int, tileWidth, tileHeight int) int {
	sX := x / tileWidth
	sY := y / tileHeight
	return (sY * width) + sX
}
