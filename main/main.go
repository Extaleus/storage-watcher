package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	mainApp := app.New()
	mainW := mainApp.NewWindow("Hello World")
	mainW.Resize(fyne.NewSize(800, 450))
	mainW.CenterOnScreen()

	// Создаем два поля для ввода
	entry1 := widget.NewEntry()
	entry1.SetPlaceHolder("input folder path")
	// entry1.Resize(fyne.NewSize(400, entry1.MinSize().Height))
	entry2 := widget.NewEntry()
	entry2.SetPlaceHolder("output folder path")
	// entry2.Resize(fyne.NewSize(400, entry1.MinSize().Height))

	// Создаем две кнопки для выбора файлов
	button1 := widget.NewButton("Select Input Folder", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				entry1.SetText(uri.Path())
			}
		}, mainW).Show()
	})

	button2 := widget.NewButton("Select Output Folder", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				entry2.SetText(uri.Path())
			}
		}, mainW).Show()
	})

	// Создаем кнопку внизу окна
	bottomButton := widget.NewButton("Process Files", func() {
		log.Println("File 1:", entry1.Text)
		log.Println("File 2:", entry2.Text)
	})

	inputFolder := container.NewHBox(
		entry1,
		button1,
	)
	inputFolder.Resize(fyne.NewSize(400, inputFolder.MinSize().Height))

	outputFolder := container.NewHBox(
		entry2,
		button2,
	)
	outputFolder.Resize(fyne.NewSize(400, outputFolder.MinSize().Height))

	// Создаем вертикальный бар с пятью кнопками
	filterButtons := container.NewVBox(
		widget.NewButton(".png", func() { fmt.Println("Filter 1 selected") }),
		widget.NewButton(".jpeg", func() { fmt.Println("Filter 3 selected") }),
		widget.NewButton(".webp", func() { fmt.Println("Filter 4 selected") }),
	)

	// Создаем контейнер для основного содержимого
	content := container.NewBorder(
		nil,           // Верх
		bottomButton,  // Низ
		nil,           // Лево
		filterButtons, // Право
		container.NewVBox( // Центр
			inputFolder,
			outputFolder,
		),
	)

	mainW.SetContent(content)
	mainW.ShowAndRun()
}
