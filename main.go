package main

import (
	"flag"
	"log"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/ncruces/zenity"
)

var AtlasPath *string
var MapImportPath *string
var TileSize *int
var TileWidth int
var TileHeight int
var ImportedLayerIndices [][]int
var IsImport *bool

func main() {
	AtlasPath = flag.String("path", "", "File path to atlas relative to the folder this program is being run from.")
	IsImport = flag.Bool("import", true, "File path to map file (for importing and editing a saved level).")
	TileSize = flag.Int("tile_size", 0, "The width and height in pixels of the tiles in the atlas (only square tiles supported).")

	flag.Parse()

	if *IsImport == true {
		filePath, _ := zenity.SelectFile(
			zenity.Title("Load Level"),
			zenity.FileFilters{
				{
					Name:     "Map Files",
					Patterns: []string{"*.map.config"},
				},
			},
		)
		MapImportPath = &filePath
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
