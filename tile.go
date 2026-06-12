package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	x, y       int
	img        *ebiten.Image
	AtlasIndex int // negative means tile does not reference an atlas
}

func NewTile(x, y int, img *ebiten.Image, atlasIndex int) *Tile {
	return &Tile{x, y, img, atlasIndex}
}

func NewEmptyTile(x, y int) *Tile {
	return &Tile{x: x, y: y, img: ebiten.NewImage(TileWidth, TileHeight), AtlasIndex: -1}
}

func (t *Tile) SetImg(img *ebiten.Image, atlasIndex int) {
	t.img = img
	t.AtlasIndex = atlasIndex
}

func (t *Tile) Reset() {
	p := t.img.Bounds().Size()
	t.img = ebiten.NewImage(p.X, p.Y)
	t.AtlasIndex = -1
}

func (t *Tile) Draw(screen *ebiten.Image, cam *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(t.x), float64(t.y))
	op.GeoM.Concat(cam.GeoM)
	screen.DrawImage(t.img, op)
}
