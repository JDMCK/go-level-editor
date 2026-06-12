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
	width, height int
	x, y          int
	tiles         []*Tile
}

func NewEmptyContainer(x, y int, width, height int) *Container {
	tiles := make([]*Tile, 0, width*height)
	for i := range width * height {
		x, y := tilePositionFromIndex(i, width)
		tiles = append(tiles, NewEmptyTile(x, y))
	}
	return &Container{
		width:  width,
		height: height,
		x:      x,
		y:      y,
		tiles:  tiles,
	}
}

func NewContainerFromAtlas(x, y int, path string) *Container {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal("Failed to read in tile set.")
	}

	size := img.Bounds().Size()
	width := size.X / TileWidth
	height := size.Y / TileHeight

	c := Container{
		height: height,
		width:  width,
		tiles:  make([]*Tile, 0, width*height),
		x:      x,
		y:      y,
	}

	for i := range width * height { // loop through width * height to fill out a full rectangle
		rx, ry := tilePositionFromIndex(i, width)
		rect := image.Rect(rx, ry, rx+TileWidth, ry+TileHeight)
		img := img.SubImage(rect).(*ebiten.Image)
		c.tiles = append(c.tiles, NewTile(rx+x, ry+y, img, i))
	}

	return &c
}

func (c *Container) TileFromCursor(op *ebiten.DrawImageOptions) *Tile {
	newOp := c.worldOp()
	newOp.GeoM.Concat(op.GeoM)
	cx, cy := CursorPosition(newOp)

	index := TileIndexFromPosition(cx, cy, c.width, c.height)
	if index < 0 || index >= c.width*c.height {
		return nil
	}
	tile := c.tiles[index]
	return tile
}

// func (c *Container) TileFromPosition(x, y int, op *ebiten.DrawImageOptions) *Tile {
// 	newOp := c.worldOp()
// 	newOp.GeoM.Concat(op.GeoM)

// 	// Invert so Apply maps screen-space back into local tile space.
// 	newOp.GeoM.Invert()
// 	tx, ty := newOp.GeoM.Apply(float64(x), float64(y))

// 	count := c.width * c.height
// 	index := TileIndexFromPosition(int(tx), int(ty), c.width, count/c.width)
// 	if index < 0 || index >= count {
// 		return nil
// 	}
// 	tile := c.tiles[index]
// 	return tile
// }

func tilePositionFromIndex(index int, width int) (int, int) {
	return (index % width) * TileWidth, (index / width) * TileHeight
}

func (c *Container) worldOp() *ebiten.DrawImageOptions {
	worldOp := ebiten.DrawImageOptions{}
	worldOp.GeoM.Translate(float64(c.x), float64(c.y))
	return &worldOp
}

func (c *Container) Clear() {
	for _, t := range c.tiles {
		t.Reset()
	}
}

// changes the width of the container while preserving the tiles already there
func (c *Container) SetWidth(newWidth int) {
	d := newWidth - c.width
	if d == 0 {
		return
	}
	var newTiles []*Tile
	if d > 0 { // expanding width
		newTiles = make([]*Tile, newWidth*c.height)
		for row := range c.height {
			for col := range c.width { // copy the original columns first
				i1 := indexFromCoords(c.width, row, col)
				i2 := indexFromCoords(newWidth, row, col)
				newTiles[i2] = c.tiles[i1]
			}
			// assume size will only expand by one at a time
			i := indexFromCoords(newWidth, row, newWidth-1)
			x, y := tilePositionFromIndex(i, newWidth)
			newTiles[i] = NewEmptyTile(x, y)
		}
	}
	if d < 0 { // shrinking width
		newTiles = make([]*Tile, newWidth*c.height)
		for row := range c.height {
			for col := range newWidth {
				i1 := indexFromCoords(c.width, row, col)
				i2 := indexFromCoords(newWidth, row, col)
				newTiles[i2] = c.tiles[i1]
			}
		}
	}
	c.width = newWidth
	c.tiles = newTiles
}

// changes the height of the container while preserving the tiles already there
func (c *Container) SetHeight(newHeight int) {
	d := newHeight - c.height
	if d == 0 {
		return
	}
	if d > 0 { // expanding height
		for i := range c.width {
			x, y := tilePositionFromIndex(c.height*c.width+i, c.width)
			c.tiles = append(c.tiles, NewEmptyTile(x, y))
		}
	}
	if d < 0 { // shrinking height
		c.tiles = c.tiles[:len(c.tiles)-c.width]
	}
	c.height = newHeight
}

func indexFromCoords(width int, row, col int) int {
	return row*width + col
}

func (c *Container) Update() {}

func (c *Container) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	newOp := c.worldOp()
	newOp.GeoM.Concat(op.GeoM)

	for _, t := range c.tiles {
		t.Draw(screen, op)
	}

	w := c.width * TileWidth
	h := c.height * TileHeight
	x0, y0 := newOp.GeoM.Apply(0, 0)
	x1, y1 := newOp.GeoM.Apply(float64(w), float64(h))

	vector.StrokeRect(screen,
		float32(x0), float32(y0),
		float32(x1-x0), float32(y1-y0),
		BorderWidth, color.White, false)
}
