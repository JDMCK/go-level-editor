package main

import (
	"fmt"
	gui "level-editor/gui"
	"log"

	eb "github.com/hajimehoshi/ebiten/v2"
)

type Editor struct {
	GUI []gui.Element
}

var editor *Editor

func init() {
	guiEls := make([]gui.Element, 0)
	guiEls = append(guiEls, gui.NewBasicButton("Save", 50, 50, gui.Large, gui.Primary, func() { fmt.Println("Saved") }))
	guiEls = append(guiEls, gui.NewBasicButton("Cancel", 50, 100, gui.Small, gui.Secondary, func() { fmt.Println("Cancel") }))
	guiEls = append(guiEls, gui.NewBasicButton("Delete", 50, 150, gui.Medium, gui.Danger, func() { fmt.Println("Delete") }))
	guiEls = append(guiEls, gui.NewNumberPicker(1, 5, 0, 50, 200))
	guiEls = append(guiEls, gui.NewCheckbox(50, 250))
	editor = &Editor{
		GUI: guiEls,
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
