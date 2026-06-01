package main

import (
	"log"

	eb "github.com/hajimehoshi/ebiten/v2"
)

func main() {
	editor := NewEditor()
	eb.SetWindowSize(ScreenWidth, ScreenHeight)
	eb.SetWindowResizingMode(eb.WindowResizingModeEnabled)
	eb.SetWindowTitle("Level Editor")

	if err := eb.RunGame(editor); err != nil {
		log.Fatal(err)
	}
}
