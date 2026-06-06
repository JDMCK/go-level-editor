package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// highlights tiles
type Cursor struct {
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

func (c *Cursor) SelectTile(tile *Tile) {
	if tile == nil {
		c.enabled = false
		c.tile = tile
		return
	}
	c.enabled = true
	c.tile = tile
	c.width = 1
	c.height = 1
}

func (c *Cursor) MultiSelect(tile *Tile, width, height int) {
	if tile == nil {
		c.enabled = false
		c.tile = tile
		return
	}
	c.enabled = true
	c.tile = tile
	c.width = width
	c.height = height
}

func (c *Cursor) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if c.enabled == false || c.tile == nil {
		return
	}
	c.img.Clear()

	newOp := ebiten.DrawImageOptions{}
	newOp.GeoM.Translate(float64(c.tile.x), float64(c.tile.y))
	newOp.GeoM.Concat(op.GeoM)

	w := c.width * c.tileWidth
	h := c.height * c.tileHeight
	x0, y0 := newOp.GeoM.Apply(0, 0)
	x1, y1 := newOp.GeoM.Apply(float64(w), float64(h))

	vector.StrokeRect(screen,
		float32(x0), float32(y0),
		float32(x1-x0), float32(y1-y0),
		HighlightWidth, color.White, false)
}
