package main

import (
	"flag"
	"log"

	eb "github.com/hajimehoshi/ebiten/v2"
)

var AtlasPath *string
var TileSize *int
var TileWidth int
var TileHeight int

func main() {
	AtlasPath = flag.String("path", "", "File path to atlas relative to the folder this program is being run from.")
	TileSize = flag.Int("tile_size", 0, "The width and height in pixels of the tiles in the atlas (only square tiles supported).")

	flag.Parse()

	if *AtlasPath == "" ||
		*TileSize == 0 {
		log.Fatal("Invalid or missing params. Use -h or -help for more information.")
	}

	TileWidth = *TileSize
	TileHeight = *TileSize

	editor := NewEditor()
	eb.SetWindowSize(ScreenWidth, ScreenHeight)
	eb.SetWindowResizingMode(eb.WindowResizingModeEnabled)
	eb.SetWindowTitle("Level Editor")

	if err := eb.RunGame(editor); err != nil {
		log.Fatal(err)
	}
}
