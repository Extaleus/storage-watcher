package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// moveImageFiles перемещает все файлы с допустимыми расширениями из sourceDir в destinationDir
func MoveAllImageFiles(sourceDir, destinationDir string, typesFilterValues []bool) ([]string, error) {
	var movedFiles []string

	// Проверяем, существует ли исходная папка
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return movedFiles, fmt.Errorf("исходная папка %s не существует", sourceDir)
	}

	// Проверяем, существует ли целевая папка, и создаем её, если нет
	if _, err := os.Stat(destinationDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
			return movedFiles, fmt.Errorf("не удалось создать целевую папку %s: %v", destinationDir, err)
		}
	}

	// Читаем файлы из исходной папки
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return movedFiles, fmt.Errorf("не удалось прочитать исходную папку %s: %v", sourceDir, err)
	}

	// Перемещаем файлы
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(sourceDir, file.Name())
			if isImageFile(filePath, typesFilterValues) {
				destinationPath := filepath.Join(destinationDir, file.Name())
				if err := os.Rename(filePath, destinationPath); err != nil {
					return movedFiles, fmt.Errorf("не удалось переместить файл %s: %v", file.Name(), err)
				}
				movedFiles = append(movedFiles, "Перемещён: "+destinationDir+"/"+file.Name())
			}
		}
	}

	return movedFiles, nil
}
