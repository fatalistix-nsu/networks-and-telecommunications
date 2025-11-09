package main

import (
	"github.com/fatalistix/copies-detector/controller"

	"fyne.io/fyne/v2/app"
)

func main() {
	application := app.New()
	window := application.NewWindow("Multicaster")
	cntrlr := controller.NewController(window)
	cntrlr.Init()
	window.ShowAndRun()
}
