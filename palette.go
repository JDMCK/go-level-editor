package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const PaletteScale = 4

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

func (p *Palette) Update() {
	// update cursor (draw / erase)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(PaletteScale, PaletteScale)
	// curTile, x, y := p.container.TileFromCursor(op)
	// p.cursor.SelectTile(x, y, curTile)
	if p.cursor.tile != nil && ebiten.IsMouseButtonPressed(Primary) {
		p.cursor.GetTile().img.Fill(color.RGBA{255, 255, 0, 255})
	}
	if p.cursor.tile != nil && ebiten.IsMouseButtonPressed(Secondary) {
		p.cursor.GetTile().img.Clear()
	}
}

func (p *Palette) SelectedTile() *Tile {
	return p.cursor.tile
}

func (p *Palette) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(PaletteScale, PaletteScale)
	p.container.Draw(screen, &op)
	p.cursor.Draw(screen, &op)
}
