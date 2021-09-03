package xsd_validation

import (
	"errors"
	"fmt"
	"github.com/lestrrat-go/libxml2"
	"github.com/lestrrat-go/libxml2/xsd"
	"io"
	"io/ioutil"
	"log"
	"os"
)



func ValidateXSD(reestr *io.Reader, xsdfile *os.File, filename string) error {
	xsdbuf, err := ioutil.ReadAll(xsdfile)
	if err != nil {
		log.Printf("Не удалось открыть файл XSD: %s", err)
		return err
	}

	bufxml, err := ioutil.ReadAll(*reestr)
	if err != nil {
		log.Printf("Не удалось открыть файл реестра: %s", err)
		return err
	}
	s, err := xsd.Parse(xsdbuf)
	if err != nil {
		log.Printf("Ошибка чтения XSD: %s", err)
		return err
	}
	defer s.Free()

	d, err := libxml2.Parse(bufxml)
	if err != nil {
		return errors.New(fmt.Sprintf("Ошибка чтения XML файла %s: %s", filename, err))
	}
	defer d.Free()

	if err := s.Validate(d); err != nil {
		for _, e := range err.(xsd.SchemaValidationError).Errors() {
			log.Printf("error: %s", e.Error())
		}
		return err
	}
	return nil
}
