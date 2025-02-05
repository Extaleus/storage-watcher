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

	sourceDir := widget.NewEntry()
	sourceDir.SetText("C:/Users/extaleus/Downloads/tests")
	sourceDir.SetPlaceHolder("Путь к исходной папке")
	buttonSelectSourceDir := widget.NewButton("Выбрать папку", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				sourceDir.SetText(uri.Path())
			}
		}, mainW).Show()
	})

	destinationDir := widget.NewEntry()
	destinationDir.SetText("C:/Users/extaleus/Downloads/tests_copy")
	destinationDir.SetPlaceHolder("Путь к папке назначения")
	buttonSelectDestinationDir := widget.NewButton("Выбрать папку", func() {
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				destinationDir.SetText(uri.Path())
			}
		}, mainW).Show()
	})

	buttonSwitchDirs := widget.NewButton("Swap dirs", func() {
		tempDir := destinationDir.Text
		destinationDir.SetText(sourceDir.Text)
		sourceDir.SetText(tempDir)
	})

	filenamesFilter := widget.NewEntry()
	filenamesFilter.SetPlaceHolder("Фильтрация по имени (или части) в виде: \"лилия\", \"помидор\".")

	var logBuilder strings.Builder
	logBuilder.WriteString(" Пусто")
	logContainer := widget.NewTextGridFromString(logBuilder.String())

	// Функция для обновления logContainer
	updateLogContainer := func(logData []string) {
		logBuilder.Reset() // Очищаем builder
		for _, line := range logData {
			logBuilder.WriteString(" " + line + "\n")
		}
		if len(logData) == 0 {
			logBuilder.WriteString(" Ничего не произошло")
		}
		logContainer.SetText(logBuilder.String())
	}

	scrollContainer := container.NewScroll(logContainer)
	InputPathBox := container.NewVBox(sourceDir, destinationDir)
	OutputPathBox := container.NewVBox(buttonSelectSourceDir, buttonSelectDestinationDir)
	typesFilter := createFilterButtons(typesFilterValues)

	executionTimeLabel := widget.NewLabel("Время выполнения: 0.00 сек.")

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
				buttonSwitchDirs,
				layout.NewSpacer(),
				widget.NewButton("Удалить дубликаты", func() {
					start := time.Now()
					logData, err := utils.FindAndRemoveDuplicates(sourceDir.Text, typesFilterValues)
					elapsed := time.Since(start)
					seconds := elapsed.Seconds()
					executionTimeLabel.SetText(fmt.Sprintf("Время выполнения: %.2f сек.", seconds))

					updateLogContainer(logData)

					if err != nil {
						fmt.Printf("Error: %v\n", err)
					} else {
						fmt.Println("Duplicate files removed successfully.")
					}
				}),
				layout.NewSpacer(),
				widget.NewButton("Переместить все файлы", func() {
					start := time.Now()
					logData, err := utils.MoveAllImageFiles(sourceDir.Text, destinationDir.Text, typesFilterValues)
					elapsed := time.Since(start)
					seconds := elapsed.Seconds()
					executionTimeLabel.SetText(fmt.Sprintf("Время выполнения: %.2f сек.", seconds))

					updateLogContainer(logData)

					if err != nil {
						fmt.Printf("Error: %v\n", err)
					} else {
						fmt.Println("Файлы перемещены успешно")
					}
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
					container.NewHBox(widget.NewLabel("Отчёт о выполнении:"), layout.NewSpacer(), createColoredSeparator(), executionTimeLabel),
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
		widget.NewLabel("Выберите расширения файлов для взаимодействия:"),
		widget.NewCheck(".png", func(b bool) { typesFilterValues[0] = b }),
		widget.NewCheck(".jpeg", func(b bool) { typesFilterValues[1] = b }),
		widget.NewCheck(".webp", func(b bool) { typesFilterValues[2] = b }),
		widget.NewCheck(".svg", func(b bool) { typesFilterValues[3] = b }),
	)
}
