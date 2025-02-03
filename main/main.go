package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	utils "github.com/Extaleus/storage-watcher/utils"
)

func main() {
	mainApp := app.New()
	mainW := mainApp.NewWindow("Hello World")
	mainW.Resize(fyne.NewSize(800, 450))
	mainW.CenterOnScreen()

	// var isPng, isJpeg, isWebp, isSvg = false, false, false, false

	typesFilterValues := []bool{false, false, false, false}

	inputPath := widget.NewEntry()
	inputPath.SetText("C:/Users/extaleus/Downloads/tests")
	inputPath.SetPlaceHolder("Input directory")
	selectInputPathBtn := widget.NewButton("Select directory", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				inputPath.SetText(uri.Path())
			}
		}, mainW).Show()
	})

	outputPath := widget.NewEntry()
	outputPath.SetText("C:/Users/extaleus/Downloads/images")
	outputPath.SetPlaceHolder("Output directory")
	selectOutputPathBtn := widget.NewButton("Select directory", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				outputPath.SetText(uri.Path())
			}
		}, mainW).Show()
	})

	filenamesFilter := widget.NewEntry()
	filenamesFilter.SetPlaceHolder("Filter: \"lilies\", \"romashki\"")

	var logBuilder strings.Builder
	logBuilder.WriteString(" Empty")
	logContainer := widget.NewTextGridFromString(logBuilder.String())

	// Функция для обновления logContainer
	updateLogContainer := func(logData []string) {
		logBuilder.Reset() // Очищаем builder
		for _, line := range logData {
			logBuilder.WriteString(" " + line + "\n")
		}
		if len(logData) == 0 {
			logBuilder.WriteString(" Nothing was deleted")
		}
		logContainer.SetText(logBuilder.String())
	}

	scrollContainer := container.NewScroll(logContainer)
	InputPathBox := container.NewVBox(inputPath, outputPath)
	OutputPathBox := container.NewVBox(selectInputPathBtn, selectOutputPathBtn)
	typesFilter := createFilterButtons(typesFilterValues)

	executionTimeLabel := widget.NewLabel("Execution time: 0.00 sec.")

	pathsContent := container.NewBorder(
		nil, createColoredSeparator(), nil, OutputPathBox,
		InputPathBox,
	)

	filtersContent := container.NewVBox(
		pathsContent,
		filenamesFilter,
		typesFilter,
	)

	content := container.NewBorder(
		container.NewBorder(
			createColoredSeparator(),
			nil,
			createColoredSeparator(),
			createColoredSeparator(),
			filtersContent,
		),
		container.NewBorder(
			createColoredSeparator(),
			createColoredSeparator(),
			createColoredSeparator(),
			createColoredSeparator(),
			container.NewHBox(
				layout.NewSpacer(),
				widget.NewButton("Remove duplicates", func() {
					start := time.Now()
					logData, err := utils.FindAndRemoveDuplicates(inputPath.Text, typesFilterValues)
					elapsed := time.Since(start)
					// Форматируем время в секунды и миллисекунды
					seconds := elapsed.Seconds()

					updateLogContainer(logData)
					executionTimeLabel.SetText(fmt.Sprintf("Execution time: %.2f sec.", seconds))

					if err != nil {
						fmt.Printf("Error: %v\n", err)
					} else {
						fmt.Println("Duplicate files removed successfully.")
					}
				}),
				layout.NewSpacer(),
				widget.NewButton("Move all files", func() {
					utils.MoveAllFiles(inputPath.Text, outputPath.Text, typesFilterValues)
				}),
				layout.NewSpacer(),
			),
		),
		nil,
		nil,
		container.NewBorder(
			createColoredSeparator(),
			nil,
			createColoredSeparator(),
			createColoredSeparator(),
			container.NewBorder(
				container.NewBorder(
					nil,
					createColoredSeparator(),
					nil,
					nil,
					container.NewHBox(widget.NewLabel("Log:"), layout.NewSpacer(), createColoredSeparator(), executionTimeLabel),
				),
				nil,
				nil,
				nil,
				scrollContainer,
			),
		),
	)

	mainW.SetContent(content)
	mainW.ShowAndRun()
}

// createColoredSeparator создает цветной разделитель
func createColoredSeparator() fyne.CanvasObject {
	return canvas.NewRectangle(color.White)
}

// createFilterButtons создает компонент с кнопками фильтров
func createFilterButtons(typesFilterValues []bool) *fyne.Container {
	return container.NewHBox(
		widget.NewLabel("Select the file extensions to scan:"),
		widget.NewCheck(".png", func(b bool) { typesFilterValues[0] = b }),
		widget.NewCheck(".jpeg", func(b bool) { typesFilterValues[1] = b }),
		widget.NewCheck(".webp", func(b bool) { typesFilterValues[2] = b }),
		widget.NewCheck(".svg", func(b bool) { typesFilterValues[3] = b }),
	)
}
