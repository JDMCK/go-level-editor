package main

import (
	"log"

	eb "github.com/hajimehoshi/ebiten/v2"
)

func main() {
	editor := NewEditor()
	eb.SetWindowSize(800, 800)
	eb.SetWindowTitle("Level Editor")

	if err := eb.RunGame(editor); err != nil {
		log.Fatal(err)
	}
}
