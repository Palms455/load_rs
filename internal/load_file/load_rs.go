package load_file

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"load_rs/configs"
	"log"
	"os"
	"path/filepath"
	"load_rs/internal/xml_parse"
	"load_rs/internal/xsd_validation"
)

func LoadRS(filename string, i int) {
	// Распаковка, валидация реестров
	// Open a zip archive for reading.
	folder := configs.GetFolder()
	r, err := zip.OpenReader(filepath.Join(folder, filename))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer r.Close()

	// Iterate through the files in the archive,
	rs_data := &xml_parse.RsFile{}
	for _, f := range r.File {
		log.Printf("Обработка файла %s", f.Name)
		log.Printf("WORKER RUNING %d\n", i)

		reestr_file, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer reestr_file.Close()

		xsd_name := string([]rune(f.Name)[0]) + ".xsd"
		xsdfile := filepath.Join("configs/xsd", xsd_name)
		xsd, err := os.Open(xsdfile)
		if err != nil {
			log.Printf("Не удалось открыть файл: %s", err)
			return
		}
		defer xsd.Close()
		var buf bytes.Buffer
		tee := io.TeeReader(reestr_file, &buf)

		if err := xsd_validation.ValidateXSD(&tee, xsd); err != nil {
			log.Printf("error: %s", err.Error())
		}

		rs_data.XmlDecode(&buf, string([]rune(f.Name)[0]))
		//json_data, _ := rs_data.MarshalJSON()
		json_data, _ := json.Marshal(rs_data)
		fmt.Sprintf("%s", string(json_data))

	}
	log.Printf("Файл %s обработан!", filename)
	return
}
