package load_file

import (
	"archive/zip"
	"bytes"
	"github.com/lestrrat-go/libxml2/xsd"
	"io"
	"load_rs/internal/storage"
	"load_rs/internal/xml_parse"
	"load_rs/internal/xsd_validation"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func LoadRS(file string, i int) {
	// Распаковка, валидация реестров
	// Open a zip archive for reading.
	r, err := zip.OpenReader(file)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer os.Remove(file)
	defer r.Close()

	rs := &storage.RsFile{}
	rs.GetFromFile(file)

	open_period, err := storage.GetCurrentPeriod(rs.Period)
	if open_period == 0 {
		log.Fatalf("Период %s файла %s закрыт для загрузки или не существует", rs.Period, rs.Filename)
		return
	}

	errArray := []string{}

	// Iterate through the files in the archive,
	rs_data := &xml_parse.RsFile{}
	for _, f := range r.File {
		log.Printf("Обработка файла %s", f.Name)

		reestr_file, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer reestr_file.Close()

		xsd_name := string([]rune(f.Name)[0]) + ".xsd"
		xsd_path := filepath.Join("configs/xsd", xsd_name)
		xsd_file, err := os.Open(xsd_path)
		if err != nil {
			log.Printf("Не удалось открыть файл: %s", err)
			return
		}
		defer xsd_file.Close()
		var buf bytes.Buffer
		tee := io.TeeReader(reestr_file, &buf)

		if err := xsd_validation.ValidateXSD(&tee, xsd_file); err != nil {
			for _, e := range err.(xsd.SchemaValidationError).Errors() {
				errArray = append(errArray, string(e.Error()))
				log.Printf("error: %s", e.Error())
			}
		}
		if len(errArray) == 0 {
			rs_data.XmlDecode(&buf, string([]rune(f.Name)[0]))
		}
	}

	rs.ErrorMsg = map[string]string{
		"load_errors": strings.Join(errArray, ". "),
	}
	if len(errArray) == 0 {
		rs.Status = "20"
	} else {
		rs.Status = "-1"
	}

	dst := filepath.Join("uploads", rs.Rsfile)
	if err := copy_file(file, dst); err != nil {
		log.Fatal(err)
		return
	}
	rs.Rs_data = *rs_data
	if err := rs.LoadRs(); err != nil {
		return
	}
	log.Printf("Файл %s обработан!", file)
	return
}
