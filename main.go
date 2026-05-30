package main

import (
	"fmt"
	"level-editor/ui"
	"log"

	eb "github.com/hajimehoshi/ebiten/v2"
)

type Editor struct {
	GUI []ui.Element
}

var editor *Editor

func init() {
	gui := make([]ui.Element, 0)
	gui = append(gui, ui.NewBasicButton("Save", 50, 50, ui.Large, ui.Primary, func() { fmt.Println("Saved") }))
	gui = append(gui, ui.NewBasicButton("Cancel", 50, 100, ui.Small, ui.Secondary, func() { fmt.Println("Cancel") }))
	gui = append(gui, ui.NewBasicButton("Delete", 50, 150, ui.Medium, ui.Danger, func() { fmt.Println("Delete") }))
	editor = &Editor{
		GUI: gui,
	}
}

func (e *Editor) Update() error {
	for _, el := range e.GUI {
		el.Update()
	}
	return nil
}

func (e *Editor) Draw(screen *eb.Image) {
	for _, el := range e.GUI {
		el.Draw(screen)
	}
}

func (e *Editor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800
}

func main() {
	eb.SetWindowSize(800, 800)
	eb.SetWindowTitle("Level Editor")

	if err := eb.RunGame(editor); err != nil {
		log.Fatal(err)
	}
}
