package utils

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

func FindAndRemoveDuplicates(dirPath string, typesFilterValues []bool) ([]string, error) {
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
		if !info.IsDir() && isImageFile(path, typesFilterValues) {
			fileChan <- path
		}
		return nil
	})

	close(fileChan) // Закрываем канал после завершения обхода файлов
	wg.Wait()       // Ждем завершения всех горутин

	if err != nil {
		return nil, err
	}

	numberOfFiles := 0
	for _, files := range hashGroups {
		if len(files) > 1 {
			numberOfFiles += len(files) - 1
		}
	}
	fmt.Printf("Общее количество элементов: %d\n", numberOfFiles)

	removedFilesNamesList := make([]string, numberOfFiles)

	j := 0
	// Удаляем дубликаты
	for _, files := range hashGroups {
		if len(files) > 1 {
			// Сортируем файлы для детерминированного удаления
			sort.Strings(files)
			// Оставляем первый файл, остальные удаляем
			for i := 0; i < len(files)-1; i++ {
				if _, err := os.Stat(files[i]); err == nil {
					removedFilesNamesList[j] = "Removed: " + files[i]
					j++
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

	return removedFilesNamesList, nil
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
func isImageFile(path string, typesFilterValues []bool) bool {
	ext := strings.ToLower(filepath.Ext(path))
	allowedExtensions := createAllowedExtensions(typesFilterValues)
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}

	return false
}

func createAllowedExtensions(typesFilterValues []bool) []string {
	allowedExtensions := []string{}

	// Порядок расширений: .png, .jpeg, .webp, .svg
	if len(typesFilterValues) > 0 && typesFilterValues[0] {
		allowedExtensions = append(allowedExtensions, ".png")
	}
	if len(typesFilterValues) > 1 && typesFilterValues[1] {
		allowedExtensions = append(allowedExtensions, ".jpeg", ".jpg") // .jpeg и .jpg считаются одинаковыми
	}
	if len(typesFilterValues) > 2 && typesFilterValues[2] {
		allowedExtensions = append(allowedExtensions, ".webp")
	}
	if len(typesFilterValues) > 3 && typesFilterValues[3] {
		allowedExtensions = append(allowedExtensions, ".svg")
	}

	return allowedExtensions
}
