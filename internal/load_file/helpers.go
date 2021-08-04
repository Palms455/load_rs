package load_file

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var RgxPath = regexp.MustCompile(`(-\d\.(ZIP|zip))`)

func create_upload_dir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	return
}

func copy_file(src, dst string) (err error) {
	create_upload_dir(dst)

	in, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		log.Printf("Ошибка копирования файла %s, %s ", dst, err)
		return
	}
	err = out.Sync()
	return
}
