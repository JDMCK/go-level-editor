package main

import "github.com/hajimehoshi/ebiten/v2"

type Tile struct {
	img *ebiten.Image
}

func (t *Tile) Draw(x, y int, screen *ebiten.Image, cam *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.GeoM.Concat(cam.GeoM)
	screen.DrawImage(t.img, op)
}
