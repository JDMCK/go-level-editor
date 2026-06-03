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
	return &Container{
		count:      width * height,
		width:      width,
		tileWidth:  tileWidth,
		tileHeight: tileHeight,
		x:          x,
		y:          y,
		tiles:      make([]*Tile, 0, width*height),
	}
}

func NewContainerFromAtlas(x, y int, atlasPath string, tileCount, rowSize, tileWidth, tileHeight int) *Container {
	img, _, err := ebitenutil.NewImageFromFile(atlasPath)
	if err != nil {
		log.Fatal("Failed to read in tile set.")
	}

	c := Container{
		count: tileCount,
		width: tileWidth,
		tiles: make([]*Tile, 0, tileCount),
	}
	for i := range tileCount {
		x, y := TilePositionFromIndex(i, rowSize, tileWidth, tileHeight)
		rect := image.Rect(x, y, x+tileWidth, y+tileHeight)
		c.tiles = append(c.tiles, NewTile(img.SubImage(rect).(*ebiten.Image)))
	}

	return &c
}

func (c *Container) TileFromCursor(cam *Camera) *Tile {
	cx, cy := CursorPosition(cam)
	index := TileIndexFromPosition(cx, cy, c.width, c.tileWidth, c.tileHeight)
	return c.tiles[index]
}

func (c *Container) Update() {}

func (c *Container) Draw(screen *ebiten.Image, cam *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.x), float64(c.y))
	op.GeoM.Concat(cam.GeoM)
	height := c.tileHeight * (c.count / c.width)
	cImg := ebiten.NewImage(c.width*c.tileWidth, height)
	cImg.Fill(color.Black)
	borderImg := ebiten.NewImage(BorderWidth*2+c.width*c.tileWidth, BorderWidth*2+height)
	vector.StrokeRect(borderImg, float32(c.x+BorderWidth), float32(c.y+BorderWidth), float32(c.width*c.tileWidth), float32(height), BorderWidth, color.White, false)
	borderImg.DrawImage(cImg, op)
	screen.DrawImage(borderImg, op)
}
