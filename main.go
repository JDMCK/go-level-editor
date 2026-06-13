package main

import (
	"flag"
	"log"

	eb "github.com/hajimehoshi/ebiten/v2"
)

var AtlasPath *string
var MapImportPath *string
var TileSize *int
var TileWidth int
var TileHeight int
var ImportedLayerIndices [][]int

func main() {
	AtlasPath = flag.String("path", "", "File path to atlas relative to the folder this program is being run from.")
	MapImportPath = flag.String("map", "", "File path to map file (for importing and editing a saved level).")
	TileSize = flag.Int("tile_size", 0, "The width and height in pixels of the tiles in the atlas (only square tiles supported).")

	flag.Parse()

	if *MapImportPath != "" {
		ImportMap()
	} else if *AtlasPath == "" ||
		*TileSize == 0 {
		log.Fatal("Invalid or missing params. Use -h or -help for more information.")
	} else {
		TileWidth = *TileSize
		TileHeight = *TileSize
	}

	editor := NewEditor()
	eb.SetWindowSize(ScreenWidth, ScreenHeight)
	eb.SetWindowResizingMode(eb.WindowResizingModeEnabled)
	eb.SetWindowTitle("Level Editor")

	if err := eb.RunGame(editor); err != nil {
		log.Fatal(err)
	}
}
