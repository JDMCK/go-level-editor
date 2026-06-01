package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// highlights tiles
type Cursor struct {
	x, y                  int
	width, height         int
	tileWidth, tileHeight int
	tile                  *Tile
	enabled               bool
}

const HighlightWidth = 3

func NewCursor(tileWidth, tileHeight int) *Cursor {
	return &Cursor{
		width:      tileWidth,
		height:     tileHeight,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
	}
}

func (c *Cursor) SelectTile(x, y int, tile *Tile) {
	c.enabled = true
	c.x = x
	c.y = y
	c.tile = tile
}

func (c *Cursor) drawHighlight(screen *ebiten.Image, cam *ebiten.DrawImageOptions) {
	cursorImg := ebiten.NewImage(HighlightWidth*2+c.width*c.tileWidth, HighlightWidth*2+c.height*c.tileHeight)
	vector.StrokeRect(cursorImg, float32(c.x), float32(c.y), float32(c.width), float32(c.height), HighlightWidth, color.White, false)
	screen.DrawImage(cursorImg, cam)
}
