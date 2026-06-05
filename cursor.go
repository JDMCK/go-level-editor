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
	img                   *ebiten.Image
}

const HighlightWidth = 2

func NewCursor(tileWidth, tileHeight int) *Cursor {
	return &Cursor{
		width:      1,
		height:     1,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		img:        ebiten.NewImage(HighlightWidth*2+tileWidth, HighlightWidth*2+tileHeight),
	}
}

func (c *Cursor) SelectTile(x, y int, tile *Tile) {
	if tile == nil {
		c.enabled = false
		c.tile = tile
		return
	}
	c.enabled = true
	c.x = x
	c.y = y
	c.tile = tile
}

func (c *Cursor) GetTile() *Tile {
	return c.tile
}

func (c *Cursor) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if c.enabled == false {
		return
	}
	c.img.Clear()
	vector.StrokeRect(c.img, 0, 0, float32(c.width*c.tileWidth), float32(c.height*c.tileHeight), HighlightWidth, color.White, false)
	newOp := ebiten.DrawImageOptions{}
	newOp.GeoM.Translate(float64(c.x), float64(c.y))
	if op != nil {
		newOp.GeoM.Concat(op.GeoM)
	}
	screen.DrawImage(c.img, &newOp)
}
