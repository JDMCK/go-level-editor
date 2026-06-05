package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Palette struct {
	container *Container
	cursor    *Cursor
}

func NewPaletteFromTileMap(x, y int, path string, tileWidth, tileHeight int, width, count int) *Palette {
	con := NewContainerFromAtlas(x, y, path, tileWidth, tileHeight, width, count)
	cur := NewCursor(tileWidth, tileHeight)
	p := Palette{
		container: con,
		cursor:    cur,
	}
	return &p
}

func (p *Palette) SelectedTile() *Tile {
	return p.cursor.tile
}

func (p *Palette) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(5, 5)
	p.container.Draw(screen, &op)
	p.cursor.Draw(screen, &op)
}
