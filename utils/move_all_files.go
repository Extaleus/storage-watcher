package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MoveAllFiles(inputPath, outputPath string, typesFilterValues []bool) {
	// Читаем содержимое папки
	files, err := os.ReadDir(inputPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Перемещаем файлы
	for _, file := range files {
		// Проверяем, что это файл (а не папка)
		if !file.IsDir() {
			// Проверяем расширение файла
			isImageFile(filepath)
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
				// Формируем полный путь к файлу
				filePath := filepath.Join(outputPath, file.Name())

				// Проверяем доступность файла
				if isFileReady(filePath) {
					moveFile(filePath, outputPath)
				} else {
					fmt.Println("File is not ready:", filePath)
				}
			}
		}
	}
}

// Функция для проверки доступности файла
func isFileReady(filePath string) bool {
	file, err := os.OpenFile(filePath, os.O_WRONLY, 0666)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// Функция для перемещения файлов
func moveFile(filePath, outputPath string) {
	// Формируем новый путь для файла
	newPath := filepath.Join(outputPath, filepath.Base(filePath))

	// Пытаемся переместить файл
	if err := os.Rename(filePath, newPath); err != nil {
		fmt.Println("Error moving file:", err)
	} else {
		fmt.Println("Moved:", filePath, "->", newPath)
	}
}
