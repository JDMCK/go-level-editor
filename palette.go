package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const PaletteScale = 3

type Palette struct {
	container       *Container
	hoverCursor     *Cursor
	selectionCursor *Cursor
}

func NewPaletteFromTileMap(x, y int, path string) *Palette {
	con := NewContainerFromAtlas(x, y, path)
	hCur := NewCursor(TileWidth, TileHeight)
	sCur := NewCursor(TileWidth, TileHeight)
	sCur.SelectTile(con.tiles[0])
	p := Palette{
		container:       con,
		hoverCursor:     hCur,
		selectionCursor: sCur,
	}
	return &p
}

func (p *Palette) Update() {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(PaletteScale, PaletteScale)

	curTile := p.container.TileFromCursor(&op)
	p.hoverCursor.SelectTile(curTile)

	if p.hoverCursor.Tile != nil && ebiten.IsMouseButtonPressed(Primary) {
		tile := p.hoverCursor.Tile
		p.selectionCursor.SelectTile(tile)
	}
}

func (p *Palette) SelectedTile() *Tile {
	return p.selectionCursor.Tile
}

func (p *Palette) TileFromAtlasIndex(i int) *Tile {
	for _, t := range p.container.tiles {
		if t.AtlasIndex == i {
			return t
		}
	}
	return nil
}

func (p *Palette) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(PaletteScale, PaletteScale)

	p.container.Draw(screen, &op)
	p.hoverCursor.Draw(screen, &op)
	p.selectionCursor.Draw(screen, &op)
}
