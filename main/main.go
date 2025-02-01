package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	mainApp := app.New()
	mainW := mainApp.NewWindow("Hello World")
	mainW.Resize(fyne.NewSize(800, 450))
	mainW.CenterOnScreen()

	var isPng, isJpeg, isWebp, isSvg = false, false, false, false

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

	logData := []string{
		"First line",
	}

	var builder strings.Builder
	for _, line := range logData {
		builder.WriteString(" " + line + "\n")
	}

	logContainer := widget.NewTextGridFromString(builder.String())
	scrollContainer := container.NewScroll(logContainer)
	InputPathBox := container.NewVBox(inputPath, outputPath)
	OutputPathBox := container.NewVBox(selectInputPathBtn, selectOutputPathBtn)
	typesFilter := createFilterButtons(&isPng, &isJpeg, &isWebp, &isSvg)

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
					fmt.Println()
					fmt.Println(inputPath.Text)
					fmt.Println()
					start := time.Now()
					err := findAndRemoveDuplicates(inputPath.Text)
					elapsed := time.Since(start)
					fmt.Printf("Время выполнения: %s\n", elapsed)
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					} else {
						fmt.Println("Duplicate files removed successfully.")
					}
				}),
				layout.NewSpacer(),
				widget.NewButton("Move all files", func() {
					log.Println("Input Folder:", inputPath)
					log.Println("Output Folder:", outputPath)
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
			container.NewVBox(
				container.NewBorder(
					nil,
					createColoredSeparator(),
					nil,
					nil,
					widget.NewLabel("Output:"),
				),
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
func createFilterButtons(isPng, isJpeg, isWebp, isSvg *bool) *fyne.Container {
	return container.NewHBox(
		widget.NewLabel("Select the file extensions to scan:"),
		widget.NewCheck(".png", func(b bool) { isPng = &b }),
		widget.NewCheck(".jpeg", func(b bool) { isJpeg = &b }),
		widget.NewCheck(".webp", func(b bool) { isWebp = &b }),
		widget.NewCheck(".svg", func(b bool) { isSvg = &b }),
	)
}

func findAndRemoveDuplicates(dirPath string) error {
	maxGoroutines := runtime.NumCPU()
	fileChan := make(chan string, maxGoroutines*10)
	hashGroups := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Запускаем пул воркеров для вычисления хешей
	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range fileChan {
				hash, err := calculateFileHash(path)
				if err != nil {
					fmt.Printf("Error calculating hash for %s: %v\n", path, err)
					continue
				}

				// Блокируем доступ к map для безопасной записи
				mu.Lock()
				hashGroups[hash] = append(hashGroups[hash], path)
				mu.Unlock()
			}
		}()
	}

	// Проходим по всем файлам в указанной папке
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isImageFile(path) {
			fileChan <- path
		}
		return nil
	})

	close(fileChan) // Закрываем канал после завершения обхода файлов
	wg.Wait()       // Ждем завершения всех горутин

	if err != nil {
		return err
	}

	// Удаляем дубликаты
	for _, files := range hashGroups {
		if len(files) > 1 {
			// Сортируем файлы для детерминированного удаления
			sort.Strings(files)

			// Оставляем первый файл, остальные удаляем
			for i := 1; i < len(files); i++ {
				if _, err := os.Stat(files[i]); err == nil {
					// fmt.Printf("Removing duplicate: %s (duplicate of %s)\n", files[i], files[0])
					err := os.Remove(files[i])
					if err != nil {
						fmt.Printf("Error removing file %s: %v\n", files[i], err)
					}
				} else if os.IsNotExist(err) {
					fmt.Printf("File %s does not exist, skipping...\n", files[i])
				} else {
					fmt.Printf("Error checking file %s: %v\n", files[i], err)
				}
			}
		}
	}

	return nil
}

// calculateFileHash вычисляет SHA-256 хеш файла с буферизацией
func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	reader := bufio.NewReader(file)
	if _, err := reader.WriteTo(hash); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// isImageFile проверяет, имеет ли файл допустимое расширение
func isImageFile(path string) bool {
	ext := filepath.Ext(path)
	switch ext {
	case ".png", ".jpeg", ".jpg", ".svg", ".webp":
		return true
	default:
		return false
	}
}
