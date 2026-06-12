package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sqweek/dialog"
)

func NewSaveAction(e *Editor) func() {
	return func() {
		filePath, _ := dialog.
			File().
			Title("Save Level").
			SetStartFile("level00.config.map").
			Filter(".map").
			Save()

		buf := bytes.Buffer{}

		// atlas info
		buf.WriteString(
			fmt.Sprintf(`atlas_path=%s
tile_width=%d
tile_height=%d
map_width=%d
map_height=%d`, *AtlasPath, TileWidth, TileHeight, CanvasWidth, CanvasHeight))
		for i, l := range e.canvas {
			buf.WriteString("\n")
			buf.WriteString(generateLayerString(i, l))
		}
		os.WriteFile(filePath, buf.Bytes(), 0644)
	}
}

func generateLayerString(i int, l *Container) string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "layer_%d=", i)

	consecTiles := 0
	consecIndex := l.tiles[0].AtlasIndex

	addTileNotation := func() {
		if consecIndex == -1 {
			fmt.Fprintf(&sb, "-%d ", consecTiles)
		} else {
			fmt.Fprintf(&sb, "%d-%d ", consecIndex, consecTiles)
		}
	}

	for _, t := range l.tiles {
		if t.AtlasIndex == consecIndex {
			consecTiles++
			continue
		}
		addTileNotation()
		consecTiles = 1
		consecIndex = t.AtlasIndex // update to next tile index
	}
	addTileNotation()

	return sb.String()
}

// func ImportMap() {
// 	data, err := os.ReadFile(*ImportPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	lines := strings.Split(string(data), "\n")

// 	var (
// 		tileWidth  int
// 		tileHeight int
// 		mapWidth   int
// 		mapHeight  int
// 		img        *eb.Image
// 	)

// 	for _, line := range lines {
// 		k, v, err := parseKV(line)
// 		if err != nil {
// 			continue // likely a comment or blank line
// 		}
// 		switch k {
// 		case "atlas_path":
// 			img, _, err = ebitenutil.NewImageFromFile(v)
// 			if err != nil {
// 				return nil, err
// 			}
// 		case "tile_width":
// 			tileWidth, err = strconv.Atoi(v)
// 			if err != nil {
// 				return nil, err
// 			}
// 		case "tile_height":
// 			tileHeight, err = strconv.Atoi(v)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}

// 	size := img.Bounds().Size()
// 	rows := size.Y / frameHeight
// 	cols := size.X / frameWidth

// 	return gfx.NewAtlas(img, rows, cols, frameWidth, frameHeight), nil
// }

func parseKV(line string) (string, string, error) {
	key, value, found := strings.Cut(line, "=")
	if !found {
		return "", "", fmt.Errorf("Invalid line %s", line)
	}
	return strings.ToLower(strings.TrimSpace(key)), strings.TrimSpace(value), nil
}
