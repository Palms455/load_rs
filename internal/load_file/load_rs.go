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
)



func LoadRS(file string, i int) {
	// Распаковка, валидация реестров

	new_filepath := RgxPath.ReplaceAllString(file, ".zip")

	if err := os.Rename(file, new_filepath); err != nil {
		log.Println("Ошибка при переименовании файла: ", err)
		return
	}

	file = new_filepath

	rs := &storage.RsFile{}
	if err := rs.GetFromFile(file); err != nil {
		log.Printf("Ошибка в наименовании файла - %s: %s", file, err)
		return
	}
	//if rs.Nn == "2" {
	//	log.Printf("Загрузка исправленной части запрещена. Файл- %s", file)
	//	return
	//}

	r, err := zip.OpenReader(file)
	if err != nil {
		log.Printf("Ошибка открытия zip - %s: %s", file, err)
		return
	}

	defer r.Close()


	year, month, disallow_load_pr_nov, disallow_load_primary_rs, err := storage.GetCurrentPeriod(rs.Period)
	if year == 0 {
		log.Printf("Период %s файла %s закрыт для загрузки или не существует", rs.Period, rs.Filename)
		return
	}
	if disallow_load_pr_nov == true && rs.Nn == "2" {
		log.Printf("Период %s файла %s закрыт для загрузки исправленной части", rs.Period, rs.Filename)
		return
	}
	if disallow_load_primary_rs == true && rs.Nn != "2" {
		log.Printf("Период %s файла %s закрыт для загрузки основной части", rs.Period, rs.Filename)
		return
	}
	rs.GetRsFilePath(year, month)

	errArray := []string{}

	// Iterate through the files in the archive,
	rs_data := &xml_parse.RsData{}
	for _, f := range r.File {
		log.Printf("Обработка файла %s", f.Name)

		reestr_file, err := f.Open()
		if err != nil {
			log.Println(err)
		}
		defer reestr_file.Close()

		xsd_name := string([]rune(f.Name)[0]) + ".xsd"
		xsd_path := filepath.Join("xsd", xsd_name)
		xsd_file, err := os.Open(xsd_path)
		if err != nil {
			log.Printf("Не удалось открыть файл: %s", err)
			return
		}
		defer xsd_file.Close()
		var buf bytes.Buffer
		tee := io.TeeReader(reestr_file, &buf)

		if xsd_error := xsd_validation.ValidateXSD(&tee, xsd_file, f.Name); xsd_error != nil {
			switch type_err := xsd_error.(type) {
			case xsd.SchemaValidationError:
				for _, e := range type_err.Errors() {
					errArray = append(errArray, string(e.Error()))
					log.Printf("error: %s", e.Error())
				}
			default:
				errArray = append(errArray, string(xsd_error.Error()))
			}
		}
		if len(errArray) == 0 {
			rs_data.XmlDecode(&buf, string([]rune(f.Name)[0]))
		}
	}
	rs.ErrorMsg = map[string][]string{
		"load_errors": errArray,
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

	r.Close()
	os.Remove(file)
	log.Printf("Файл %s обработан!", file)
	return
}
