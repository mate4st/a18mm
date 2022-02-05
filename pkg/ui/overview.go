package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func InitOverview() *fyne.Container {

	blue := color.NRGBA{R: 0, G: 0, B: 180, A: 255}

	statusHaeder := container.NewGridWithColumns(2, widget.NewLabel("upper top bar"), canvas.NewRectangle(blue))
	statusHaeder2 := container.NewGridWithColumns(2, widget.NewLabel("lower top bar"), canvas.NewRectangle(blue))

	headerMenu := container.NewGridWithRows(2,
		statusHaeder,
		statusHaeder2,
	)

	var data = []string{"mod 1", "mod 2", "mod 3"}

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewCenter(widget.NewLabel(""))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {

			labelObj := o.(*fyne.Container).Objects

			labelObj[0].(*widget.Label).SetText(data[i])

		})

	blueRa := canvas.NewRectangle(color.Transparent)

	blueRa.SetMinSize(fyne.NewSize(20, 300))

	left := container.NewCenter(blueRa)

	return container.New(layout.NewBorderLayout(headerMenu, nil, left, nil), headerMenu, left, list)
}
