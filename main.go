package main

import (
	"image/color"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Container")
	green := color.NRGBA{R: 0, G: 180, B: 0, A: 255}

	text1 := canvas.NewText("Hello", green)
	text2 := canvas.NewText("There", green)
	text3 := canvas.NewText("123", green)
	text4 := canvas.NewText("456", green)
	//content := container.NewWithoutLayout(text1, text2)
	content := container.New(layout.NewBorderLayout(text1, text2, text3, text4))

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
