package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const BorderWidth = 2

// contains tiles
type Container struct {
	count, width          int
	tileWidth, tileHeight int
	x, y                  int
	tiles                 []*Tile
}

func NewEmptyContainer(x, y int, tileWidth, tileHeight int, width, height int) *Container {

	tiles := make([]*Tile, 0, width*height)
	for range width * height {
		tiles = append(tiles, NewEmptyTile(tileWidth, tileHeight))
	}

	return &Container{
		count:      width * height,
		width:      width,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		x:          x,
		y:          y,
		tiles:      tiles,
	}
}

func NewContainerFromAtlas(x, y int, path string, tileWidth, tileHeight int, width, count int) *Container {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal("Failed to read in tile set.")
	}

	c := Container{
		count:      count,
		width:      width,
		tiles:      make([]*Tile, 0, count),
		x:          x,
		y:          y,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
	}
	for i := range count {
		x, y := TilePositionFromIndex(i, width, tileWidth, tileHeight)
		rect := image.Rect(x, y, x+tileWidth, y+tileHeight)
		c.tiles = append(c.tiles, NewTile(img.SubImage(rect).(*ebiten.Image)))
	}

	return &c
}

func (c *Container) TileFromCursor(op *Camera) (*Tile, int, int) {
	cx, cy := CursorPosition(op)
	index := TileIndexFromPosition(cx, cy, c.width, c.count/c.width, c.tileWidth, c.tileHeight)
	x, y := TilePositionFromIndex(index, c.width, c.tileWidth, c.tileHeight)
	if index < 0 || index >= c.count {
		return nil, 0, 0
	}
	return c.tiles[index], x, y
}

func GetHeight(width, count int) int {
	return count / width
}

func (c *Container) Update() {}

func (c *Container) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	newOp := &ebiten.DrawImageOptions{}
	newOp.GeoM.Translate(float64(c.x), float64(c.y))
	if op != nil {
		newOp.GeoM.Concat(op.GeoM)
	}
	for i, t := range c.tiles {
		x, y := TilePositionFromIndex(i, c.width, c.tileWidth, c.tileHeight)
		t.Draw(x, y, screen, newOp)
	}

	w := c.width * c.tileWidth
	h := GetHeight(c.width, c.count) * c.tileHeight

	x0, y0 := newOp.GeoM.Apply(0, 0)
	x1, y1 := newOp.GeoM.Apply(float64(w), float64(h))

	vector.StrokeRect(screen,
		float32(x0), float32(y0),
		float32(x1-x0), float32(y1-y0),
		BorderWidth, color.White, false)

}
