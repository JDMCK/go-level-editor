package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// contains tiles
type Container struct {
	count, width int
	tiles        []*Tile
}

func NewEmptyContainer(x, y int, tileWidth, tileHeight int, width, height int) *Container {
	return &Container{
		count: width * height,
		width: width,
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
		x, y := TilePosition(i, rowSize, tileWidth, tileHeight)
		rect := image.Rect(x, y, x+tileWidth, y+tileHeight)
		c.tiles = append(c.tiles, NewTile(img.SubImage(rect).(*ebiten.Image)))
	}

	return &c
}

func (c *Container) Draw(screen *ebiten.Image, cam *ebiten.DrawImageOptions) {

}
