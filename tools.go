package main

import (
	"slices"

	eb "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func handleEdits(e *Editor, curTile *Tile) {
	if inpututil.IsMouseButtonJustPressed(Tertiary) {
		// flood fill
		l := e.canvas[e.currLayer]
		i := TileIndexFromPosition(curTile.x, curTile.y, l.width, l.height)

		// bfs
		visited := make([]bool, len(l.tiles))
		queue := make([]int, 0, len(l.tiles))
		queue = append(queue, i)
		visited[i] = true
		head := 0

		for head < len(queue) {
			nexts := getNextSteps(queue[head], l.width, l.height, curTile.AtlasIndex, l)
			head++
			for _, n := range nexts {
				if visited[n] == false {
					visited[n] = true
					queue = append(queue, n)
				}
			}
		}

		for _, i := range queue {
			pTile := e.palette.SelectedTile()
			if pTile != nil {
				t := l.TileFromIndex(i)
				t.SetImg(pTile.img, pTile.AtlasIndex)
			}
		}
	}

	// add single tile
	if eb.IsMouseButtonPressed(Primary) {
		tile := e.palette.SelectedTile()
		if tile != nil {
			curTile.SetImg(tile.img, tile.AtlasIndex)
		}
	}
	// erase single tile
	if eb.IsMouseButtonPressed(Secondary) {
		e.cursor.Tile.Reset()
	}
}

// returns up to 4 possible adjacent tiles (will not go OOB)
func getNextSteps(index int, width, height int, atlasIndex int, c *Container) []int {
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

	nexts = slices.DeleteFunc(nexts, func(i int) bool { // filter out un-like tiles
		tile := c.TileFromIndex(i)
		return tile.AtlasIndex != atlasIndex
	})

	return nexts
}
