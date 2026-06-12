package main

import (
	"slices"

	eb "github.com/hajimehoshi/ebiten/v2"
)

func handleEdits(e *Editor, curTile *Tile) {
	// add single tile
	if eb.IsMouseButtonPressed(Primary) {
		tile := e.palette.SelectedTile()
		if tile != nil && curTile != nil {
			curTile.SetImg(tile.img, tile.AtlasIndex)
		}
	}
	// erase single tile
	if eb.IsMouseButtonPressed(Secondary) {
		e.cursor.Tile.Reset()
	}

	// flood fill
	l := e.canvas[e.currLayer]
	i := TileIndexFromPosition(curTile.x, curTile.y, l.width, l.height)
	nexts := getNextSteps(i, l.width, l.height, curTile.AtlasIndex)
	nexts = slices.DeleteFunc(nexts, func(i int) bool { // is bounded by different tiles
		tile := l.TileFromIndex(i)
		return tile.AtlasIndex != curTile.AtlasIndex
	})
}

// returns up to 4 possible adjacent tiles (will not go OOB)
func getNextSteps(index int, width, height int, atlasIndex int) []int {
	nexts := make([]int, 0, 4)
	if index >= width { // not at top row
		nexts = append(nexts, index-width)
	}
	if index/width < height-1 { // not at bottom row
		nexts = append(nexts, index+width)
	}
	if index%width > 0 { // not at left side
		nexts = append(nexts, index-1)
	}
	if index%width < width-1 { // not at right side
		nexts = append(nexts, index+1)
	}

	return nexts
}
