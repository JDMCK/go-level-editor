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
	for i := range width * height {
		x, y := TilePositionFromIndex(i, width, tileWidth, tileHeight)
		tiles = append(tiles, NewEmptyTile(x, y, tileWidth, tileHeight))
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
		rx, ry := TilePositionFromIndex(i, width, tileWidth, tileHeight)
		rect := image.Rect(rx, ry, rx+tileWidth, ry+tileHeight)
		c.tiles = append(c.tiles, NewTile(rx+x, ry+y, img.SubImage(rect).(*ebiten.Image)))
	}

	return &c
}

func (c *Container) TileFromCursor(op *ebiten.DrawImageOptions) *Tile {
	newOp := ebiten.DrawImageOptions{}
	newOp.GeoM.Translate(float64(c.x), float64(c.y))
	newOp.GeoM.Concat(op.GeoM)
	cx, cy := CursorPosition(&newOp)
	index := TileIndexFromPosition(cx, cy, c.width, c.count/c.width, c.tileWidth, c.tileHeight)
	if index < 0 || index >= c.count {
		return nil
	}
	tile := c.tiles[index]
	return tile
}

func (c *Container) TileFromPosition(x, y int, op *ebiten.DrawImageOptions) *Tile {
	newOp := ebiten.DrawImageOptions{}
	newOp.GeoM.Translate(float64(c.x), float64(c.y))
	newOp.GeoM.Concat(op.GeoM)

	// Invert so Apply maps screen-space back into local tile space.
	newOp.GeoM.Invert()
	tx, ty := newOp.GeoM.Apply(float64(x), float64(y))

	index := TileIndexFromPosition(int(tx), int(ty), c.width, c.count/c.width, c.tileWidth, c.tileHeight)
	if index < 0 || index >= c.count {
		return nil
	}
	tile := c.tiles[index]
	return tile
}

func GetHeight(width, count int) int {
	return count / width
}

func (c *Container) Update() {}

func (c *Container) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	newOp := ebiten.DrawImageOptions{}
	newOp.GeoM.Translate(float64(c.x), float64(c.y))
	newOp.GeoM.Concat(op.GeoM)

	for _, t := range c.tiles {
		t.Draw(screen, op)
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
