import json
import re
import os

import xmlschema
import time
import tempfile
import zipfile


# fmt:off
xmlschema_errors_translates = tuple((re.compile(pattern) if r"\1" in repl else pattern, repl) for pattern, repl in (
    (r"The content of element ([^ ]+) is not complete", r"Содержимое тэга \1 не полное"),
    (r"Unexpected child with tag ([^ ]+) at position (\d+)", r"Лишний тэг \1 с порядковым номером \2"),
    (r"The particle ([^ ]+) occurs (\d+) times but the maximum is (\d+)", r"Тэг \1 встречается \2 раз, разрешено максимум \3"),
    (r"The particle ([^ ]+) occurs (\d+) times but the minimum is (\d+)", r"Тэг \1 встречается \2 раз, разрешено минимум \3"),
    (r"Tag ([^ ]+) expected", r"Отсутствует тэг \1"),
    (r"Tag (\([^)]+\)) expected", r"Отсутствуют тэги \1"),
    ("length has to be ", "Длинна должна быть "),
    ("value length cannot be lesser than ", "Длинна строки не может быть меньше "),
    ("value length cannot be greater than ", "Длинна строки не может быть больше "),
    ("value has to be greater or equal than ", "Значение должно быть больше или равно чем "),
    ("value has to be greater than ", "Значение должно быть больше чем "),
    ("value has to be lesser or equal than ", "Значение больше меньше или равно быть больше чем "),
    ("value has to be lesser than ", "Значение больше меньше быть больше чем "),
    ("the number of digits has to be lesser or equal ", "Кол-во знаков в числе должно быть меньше или равно "),
    ("the number of fraction digits has to be lesser ", "Кол-во знаков в числе должно быть меньше "),
    ("value must be one of ", "Значение должно быть одно из "),
    ("value doesn't match any pattern of ", "Значение не соответствует не одному образцу "),
    ("value is not true with test path ", "Значение не соответствует пути "),
    ("assertion test if false", "Тест 'если' не верный"),
    ("referenced attribute has a different fixed value ", "Данный атрибут имеет другое постоянное значение"),
    (r"default value (.+) is not compatible with attribute's type", r"Значение по умолчанию \1 не соответствует типу атрибута"),
    (r"fixed value (.+) is not compatible with attribute's type", r"Постоянное значение \1 не соответствует типу атрибута"),
    ("cannot validate against xs:NOTATION directly, only against a subtype with an enumeration facet", "Не могу проверить прямо с использованием xs:NOTATION, только с подтипом с перечислениями"),
    ("missing enumeration facet in xs:NOTATION subtype", "Отсутствует перечисление в подтипе xs:NOTATION"),
    (r"attribute ([^ ]+) has a fixed value ", r"Атрибут \1 имеет постоянное значение"),
    (r"Attribute ([^ ]+): attribute 'inheritable' value change in restriction", r"Атрибут \1: изменение не наследуемого значение в ограничениях"),
    ("missing required attribute ", "Отсутствует обязательный атрибут "),
    (" is not an attribute of the XSI namespace", " не является атрибутом пространства имен XSI "),
    (" attribute not allowed for element", " атрибут не разрешен для этого тэга"),
    (r"use of attribute ([^ ]+) is prohibited", r"Использование атрибута \1 запрещено"),
    ("schemaLocation declaration after namespace start", "Определение schemaLocation после начала пространства имен"),
    ("dynamic loaded schema change the assessment", ""),
    ("cannot use an abstract element for validation", "Нельзя использовать абстрактный тэг для валидации"),
    (r"usage of ([^ ]+) is blocked", r"Использование \1 заблокировано"),
    ("element is not nillable", "Тэг не может иметь нулевое значение"),
    ("xsi:nil attribute must have a boolean value", "Атрибут xsi:nil не может иметь нулевое значение"),
    ("xsi:nil='true' but the element has a fixed value.", "Xsi:nil='true', но тэг имеет постоянное значение"),
    ("xsi:nil='true' but the element is not empty.", "Xsi:nil='true', но тэг не пустой"),
    ("character data is not allowed because content is empty", "Строковые данные не разрешены, т.к. пустое содержимое"),
    ("must have the fixed value ", "Должно иметь постоянное значение "),
    ("a simple content element can't have child elements", "Простой тэг не может иметь дочерние тэги"),
    ("must have the fixed value ", "Должен иметь постоянное значение "),
    (r"usage of ([^ ]+) with type ([^ ]+) is blocked by head element", r"Использование \1 с типом \1 заблокировано главным тэгом"),
    (r" that matches %r is not consistent with local declaration ", r", которое равно \1 не соответствует локальному определению "),
    ("an empty 'choice' group with minOccurs > 0 cannot validate any content", "Пустая группа 'choice' с minOccurs > 0 не сможет провести валидацию"),
    ("character data between child elements not allowed", "Строковые данные не разрешены между дочерними тэгами"),
    (r"XML data depth exceeded \(MAX_XML_DEPTH=([^)]+)\)", r"Глубина XML превысила значение (MAX_XML_DEPTH=\1)"),
    (" does not match any declared element of the model group", " не соответствует ни одному объявленному тэгу модели группы"),
    (" has an unknown prefix ", " имеет неизвестный префикс "),
    ("wrong content type ", "Неправильный тип содержимого "),
    (" is not an element of the schema", " не является тэгом схемы"),
    (r"the path ([^ ]+) doesn't match any element of the schema!", r"Путь \1 не соответствует не одному тэгу схемы"),
    ("unable to select an element for decoding data, provide a valid 'path' argument", "Не смог выбрать тэг для декодирования данных, введите правильный аргумент 'path'"),
    ("one or more facets are not applicable, admitted set is %r:", ""),
    (r"a (.+) or (.+) object required", r"Объект \1 или \2 обязательны"),
    ("value is not an instance of ", "Значение не является экземпляром "),
    ("invalid value ", "Неверное значение "),
    (r"unmapped prefix ([^ ]+) on QName", r"Неиспользованные префикс \1 в QName"),
    ("Duplicated xs:ID value ", "Повторяющееся значение xs:ID "),
    ("No more than one attribute of type ID should be present in an element", "Разрешен только один атрибут c типом ID в тэге"),
    (r"boolean value (.+) requires a (.+) decoder.", r"Булево значение \1 обязательно для декодера \2"),
    (r" is not an instance of ", " не является экземпляром "),
    ("Invalid value ", "Неверное значение"),
    ("no type suitable for decoding the values ", "Нет подходящего типа для декодирования "),
    ("no type suitable for encoding the object", "Нет подходящего типа для декодирования"),
    (" is not allowed here", " не разрешен здесь"),
    (r"element ([^ ]+) not found", r"тэг \1 не найден"),
    ("unavailable namespace ", "Пространство имен недоступно для "),
    (r"element ([^ ])+ is not allowed here", r"тэг \1 не разрешен здесь"),
    (r"element ([^ ])+ not found", r"тэг \1 не найден"),
    (r"attribute ([^ ])+ not allowed", r"атрибут \1 не разрешен"),
    (r"attribute ([^ ])+ not found", r"атрибут \1 не найден"),
))

syntax_errors_translates = (
    ("junk after document element", "Ошибочный тэг после заголовочного тэга"),
    ("not well-formed (invalid token)", "Ошибка форматирования (неправильный токен)"),
    ("duplicate attribute", "Повторяющийся атрибут"),
    ("mismatched tag", "Не совпадающий тэг"),
    ("unbound prefix", "Непривязанный префикс"),
)
# fmt: on


class XmlJson:
    """Класс для преобразования xml в json и обратно на основе xsd"""

    xsd = None
    error = None

    def __init__(self, xsd_source):
        """
        Конструктор класса

        :param xsd_source: URI ссылающийся на ресурс со схемой либо путь к файлу либо файловый дескриптор \
        либо строка содержащая схему либо тип Element/ElementTree.
        :type xsd_source: Element или ElementTree или str или file-like object
        :return: экземпляр класса.
        """
        self.xsd = xmlschema.XMLSchema11(xsd_source, validation="strict")
        self.error = None

    def validate(self, xml_source, path=None):
        """
        Проверяет соответствие xml-документа схеме xsd

        :param xml_source: XML данные. Может быть экземпряром класса :class:`xmlschema.XMLResource`, \
        путь к файлу или URI к данным либо файловый дескриптор либо тип Element/ElementTree \
        либо строка содержащая XML-данные.
        :param path: опциональное выражение XPath соответствующее элементам в схеме, если source_xml \
        содержит не весь документ схемы, а его часть. Если не указан, используется корневой элемент схемы.

        :raises: :exc:`xmlschema.XMLSchemaValidationError` исключение возникает если xml_source не проходит проверку xsd.
        """
        self.xsd.validate(xml_source, path=path)

    def is_valid(self, xml_source, path=None):
        """
        То же что :meth:`validate` но не вызывает исключений, а возвращает ``True`` при успешной валидации \
        и ``False`` если нет. В случае ошибки валидации, текст ошибки записывается в self.error
        """
        errors = []
        try:
            for error in self.xsd.iter_errors(xml_source, path):
                error = f"Ошибка проверки XML файла: {error.path}: {self.translate_xmlschema_errors(error.reason)}"
                errors.append(error)
                break
        except Exception as e:

            self.error = f"Синтаксическая ошибка в XML файле: {self.translate_syntax_errors(e)}"
        else:
            self.error = "\n".join(errors)
        return not self.error

    def xml_to_json(self, xml_source, path=None):
        """
        Преобразует xml в объект

        :param xml_source: XML данные. Может быть экземпряром класса :class:`xmlschema.XMLResource`, \
        путь к файлу или URI к данным либо файловый дескриптор либо тип Element/ElementTree \
        либо строка содержащая XML-данные.
        :param path: опциональное выражение XPath соответствующее элементам в схеме, если source_xml \
        содержит не весь документ схемы, а его часть. Если не указан, используется корневой элемент схемы.

        :return: List, tuple или None если валидация не прошла либо поиск по path не вернул элементов.
        """
        if self.is_valid(xml_source, path=path):
            return self.xsd.decode(xml_source, path=path, decimal_type=str)
        return None

    def xml_to_json_s(self, xml_source, path=None, indent=4):
        """
        То же что :meth:`xml_to_json` но возвращает строку
        :param indent: отступ при форматировании json
        """
        j = self.xml_to_json(xml_source, path)
        if j is not None:
            return json.dumps(j, indent=indent)
        return None

    def json_to_xml(self, data, path=None):
        """
        Преобразует json в xml

        :param data: JSON данные. List, tuple или строка.
        :param path: опциональное выражение XPath соответствующее элементам в схеме, если data \
        содержит не весь документ схемы, а его часть. Если не указан, используется корневой элемент схемы.

        :return: xml.Element.
        """
        if isinstance(data, str):
            d = json.loads(data)
        else:
            d = data
        return self.xsd.encode(d, path=path, unordered=True)

    def json_to_xml_s(self, data, path=None):
        """
        То же что :meth:`json_to_xml` но возвращает строку
        """
        return xmlschema.etree_tostring(self.json_to_xml(data, path))

    @staticmethod
    def translate_xmlschema_errors(s):
        """
        Переводит ошибки разбора xmlschema на русский язык.

        :param s: Строка ошибки.
        """
        chunks = []
        for chunk in str(s).rstrip(".").split(".", maxsplit=1):
            for translate, repl in xmlschema_errors_translates:
                try:
                    replaced, n = translate.subn(repl, chunk.strip())
                    if n:
                        chunks.append(replaced)
                        break
                except AttributeError:
                    replaced = chunk.replace(translate, repl)
                    if replaced != chunk:
                        chunks.append(replaced)
                        break
            else:
                chunks.append(chunk)
            # Проверяем это правило отдельно, чтобы не разбивать строку дополнительно:
            #   (r'attribute ([^=]+)=([^:]+): ', r"атрибут \1=\2: "),
            if chunk.startswith("attribute "):
                chunk = chunk.replace("attribute ", "атрибут ")
        return ". ".join(chunks) + "." if chunks else s

    @staticmethod
    def translate_syntax_errors(s):
        """
        Переводит синтаксические ошибки разбора на русский язык.

        :param s: Строка ошибки.
        """
        s = str(s)
        for translate, repl in syntax_errors_translates:
            replaced = s.replace(translate, repl)
            if replaced != s:
                break
        else:
            replaced = s
        replaced = replaced.replace("line", "строка")
        replaced = replaced.replace("column", "cтолбец")
        return replaced




from multiprocessing import Pool

def validate_one(filename: str) -> None:
    print(filename)
    with tempfile.TemporaryDirectory() as temp_dir:
        with zipfile.ZipFile(f"reestr/{filename}", "r") as zip:
            xml_filename = f"{filename[:-4]}.XML"
            l_file = f"L{xml_filename[1 :]}"
            xsd_file = f"{filename[0]}.xsd"
            try :
                for zipinfo in zip.infolist() :
                    zipinfo.filename = zipinfo.filename.upper()
                    zip.extract(zipinfo, temp_dir)
            except Exception as err :
                print(err)
            s = XmlJson(f"configs/xsd/{xsd_file}")
            j = s.xml_to_json(
                os.path.join(temp_dir, xml_filename),
                path="/ZL_LIST",
            )
            print(s.error)
            l = XmlJson(f"xsd/L.xsd")
            lj = s.xml_to_json(
                os.path.join(temp_dir, xml_filename)
            )
            print(l.error)


    return


if __name__ == "__main__":
    reestr = os.listdir("reestr")
    print(f"Валидация {len(reestr)} файлов")
    start = time.monotonic()
    with Pool() as p:
        p.map(validate_one, reestr)
    # s = XmlJson("test/xsd/H.xsd")
    # print(s.xsd)
    # j = s.is_valid(
    #    "test/HM220023T22_21030.XML",
    #    path="/ZL_LIST",
    # )
    # if j:
    #     print(j)
    #     xm = s.json_to_xml_s(j, path="/ZL_LIST")
    #     print(xm)
    # else:
    #print(s.error)
    end = time.monotonic()
    all_time = end-start
    print("ALL TIME: ", all_time)