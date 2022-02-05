package main

import (
	"a18mm/pkg/manager"
	"a18mm/pkg/mod"
)

func main() {

	conf, err := manager.Load()
	if err != nil {
		panic(err)
	}

	println(conf.InstallLocation)

	version, err := mod.VerifyModLoader(conf.InstallLocation, conf.LastInstalledLoaderVersion)
	if err != nil {
		panic(err)
	}

	conf.LastInstalledLoaderVersion = version
	err = conf.Save()
	if err != nil {
		panic(err)
	}

	/*
		err = mod.VerifyModLoader()

		if err != nil {
			panic(err)
		}
	*/

	/*

		application := app.NewWithID("dev.fanya.a18mm")
		window := application.NewWindow("Anno 1800 Moad settings")

		overview := ui.InitOverview()
		window.SetContent(overview)

		window.Resize(fyne.NewSize(600, 400))
		window.ShowAndRun()
	*/
}
