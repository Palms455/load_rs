package load_file

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"github.com/lestrrat-go/libxml2/xsd"
	"load_rs/internal/storage"
	"load_rs/internal/xml_parse"
	"load_rs/internal/xsd_validation"
	"log"
	"os"
	"path/filepath"
	"regexp"


)
var rgx = regexp.MustCompile(`(C|H|T|DP|DV|DO|DS|DU|DF)M(\d{6})T(\d{2})_([1-9]\d{3})(\d{1,10})\.(ZIP|zip)$`)



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

	rs := rgx.FindStringSubmatch(file)
	if len(rs) == 0 {
		log.Fatalf("Файл %s не подходит для обработки", r.File)
	}
	filename, ftype, mo, tfoms, period, nn:= rs[0], rs[1], rs[2], rs[3], rs[4], rs[5]
	fmt.Println(ftype, mo, tfoms, period, nn)

	open_period, err := storage.GetCurrentPeriod(period)
	if open_period == 0 {
		log.Fatalf("Период %s файла %s закрыт для загрузки или не существует", period, filename)
		return
	}

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
				log.Printf("error: %s", e.Error())
				return
			}
		}

		rs_data.XmlDecode(&buf, string([]rune(f.Name)[0]))
		//json_data, _ := rs_data.MarshalJSON()
		json_data, _ := json.Marshal(rs_data)


		fmt.Sprintf("%s", string(json_data))

	}
	log.Printf("Файл %s обработан!", file)
	return
}
