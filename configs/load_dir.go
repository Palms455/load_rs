package configs

import (
	"os"
	"path/filepath"
	"log"
)
// Current working directory
func GetFolder() string {
	dir, _ := os.Getwd()
	log.Println("CWD:", dir)
	return filepath.Join(dir, "reestr")
}