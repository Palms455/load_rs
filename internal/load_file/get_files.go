package load_file

import (
	"os"
	"path/filepath"
	"log"
)

func GetFolder() string {
	dir, _ := os.Getwd()
	log.Println("CWD:", dir)
	return filepath.Join(dir, "reestr")
}

func GetFiles(folder string) ([]string, error) {
	// Получение списка файлов реестров
	var files []string

	f, err := os.Open(folder)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}
