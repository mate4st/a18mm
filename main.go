package main

import (
	"fyne.io/fyne/v2/app"

	"a18mm/pkg/mod"
	"a18mm/pkg/ui"
)

func main() {

	application := app.NewWithID("dev.fanya.a18mm")
	window := application.NewWindow("Anno 1800 Moad Manager")

	overview := ui.InitOverview()
	window.SetContent(overview)

	err := mod.VerifyModLoader()

	if err != nil {
		panic(err)
	}

	// window.Resize(fyne.NewSize(600, 400))
	// window.ShowAndRun()

}
