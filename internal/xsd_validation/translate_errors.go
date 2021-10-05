package xsd_validation

import (
	"regexp"
)


// TranslateXsd Структура содержит скомпилированое регулярное выражение и строку для замены
type TranslateXsd struct {
	ErrReg     *regexp.Regexp
	ReplaceMsg string
}

func NewTranslate(errMsg, newMsg string) *TranslateXsd {
	return &TranslateXsd{
		ErrReg: regexp.MustCompile(errMsg),
		ReplaceMsg: newMsg,
	}
}

var ArrTranslate = []*TranslateXsd{
	NewTranslate(`Expected is (one of )?`, "Ожидаются теги: "),
	NewTranslate(`This element is not expected`, "Данный элемент не ожидается"),
	NewTranslate(`Missing child element\(s\)`, "Пропущен дочерний элемент"),
	NewTranslate(`Element '(\w*)': `, "Ошибка в теге '${1}': "),
	NewTranslate(`Duplicate key-sequence \[(...*)\] in unique identity-constraint '(...*)'`, "Нарушение уникальности элемента. Повторяется значение тега ${2}: ${1}"),
	NewTranslate(`unique(\w+)`, "${1}"),
	NewTranslate(`is not a valid value of the atomic type `, "- значение элемента не подходит для установленного для данного тега типа: "),
	NewTranslate(`xs:date`, "дата"),
	NewTranslate(`money-type`, "число с двумя знаками после запятой"),
	NewTranslate(`snils-type`, "СНИЛС"),
	NewTranslate(`xs:string`, "строка"),
	NewTranslate(`integer-type`, "целое число"),
	NewTranslate(`\[facet 'maxLength'\] The value has a length of ('\d+'); this exceeds the allowed maximum length of ('\d+')`, "Превышена допустимая длина (${2}) для данного тега ${1}"),
	NewTranslate(`\[facet 'minLength'\] The value has a length of ('\d+'); this underruns the allowed minimum length of ('\d+')`, "Длина элемента (${1}) менее минимально допустимой (${2})"),
	NewTranslate(`failed ...* EndTag: '</' not found`, "ошибка в структуре xml. Закрывающий тег '</' отсутствует"),
	NewTranslate(`Warning: (...*) strange happened`, "Ошибка при валидации значения элемента"),
}


func (t *TranslateXsd) ReplaceString(msg string) string {
	return t.ErrReg.ReplaceAllString(msg, t.ReplaceMsg)
}

// TranslateFlcErrors Перевод ошибок xsd внутри массива
func  TranslateFlcErrors(FlcArr []string) {
	for idx, _ := range FlcArr {
		for _, repl := range ArrTranslate {
			FlcArr[idx] = repl.ReplaceString(FlcArr[idx])
		}

	}
}

