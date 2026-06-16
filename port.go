package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ncruces/zenity"
)

func NewSaveAction(e *Editor) func() {
	return func() {
		filePath, _ := zenity.SelectFileSave(
			zenity.Title("Save Level"),
			zenity.Filename("level00.map.config"),
			zenity.FileFilters{
				{
					Name:     "Map Files",
					Patterns: []string{"*.map"},
				},
			},
		)

		buf := bytes.Buffer{}

		// atlas info
		fmt.Fprintf(&buf, `atlas_path=%s
tile_width=%d
tile_height=%d
map_width=%d
map_height=%d`, *AtlasPath, TileWidth, TileHeight, CanvasWidth, CanvasHeight)
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

func ImportMap() {
	data, err := os.ReadFile(*MapImportPath)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.SplitSeq(string(data), "\n")

	ImportedLayerIndices = make([][]int, 0, DefaultLayerCount)

	for line := range lines {
		k, v, found := parseKV(line)
		if found == false {
			continue // likely a comment or blank line
		}
		switch k {
		case "atlas_path":
			AtlasPath = &v
		case "tile_width":
			TileWidth, _ = strconv.Atoi(v)
		case "tile_height":
			TileHeight, _ = strconv.Atoi(v)
		case "map_width":
			CanvasWidth, _ = strconv.Atoi(v)
		case "map_height":
			CanvasHeight, _ = strconv.Atoi(v)
		}

		if strings.HasPrefix(k, "layer_") {
			var layer int
			fmt.Sscanf(k, "layer_%d", &layer)
			ImportedLayerIndices = append(ImportedLayerIndices, parseLayer(v))
		}
	}
}

func parseLayer(data string) []int {
	parts := strings.Split(data, " ")
	finalLen := CanvasWidth * CanvasHeight
	indices := make([]int, 0, finalLen)
	for _, p := range parts {
		atlasIndex, count, _ := strings.Cut(p, "-")
		if count == "" {
			log.Fatal("Failed to parse layer.")
		}
		countN, _ := strconv.Atoi(count)
		if atlasIndex == "" {
			for range countN {
				indices = append(indices, -1)
			}
			continue
		}
		for range countN {
			i, _ := strconv.Atoi(atlasIndex)
			indices = append(indices, i)
		}
	}
	if len(indices) != finalLen {
		log.Fatal("Failed to parse layer.")
	}
	return indices
}

func parseKV(line string) (string, string, bool) {
	key, value, found := strings.Cut(line, "=")
	return strings.ToLower(strings.TrimSpace(key)), strings.TrimSpace(value), found
}
