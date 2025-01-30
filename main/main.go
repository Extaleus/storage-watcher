package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	mainApp := app.New()
	mainW := mainApp.NewWindow("Hello World")

	mainW.SetContent(widget.NewLabel("Hello World!"))
	mainW.ShowAndRun()
}
