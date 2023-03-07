package gojson

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/Grivn/gojson/regex"
)

func JSONSchemaCamel2Snake(data []byte) []byte {
	return regex.JSONKey.ReplaceAllFunc(data, convertBytesCamel2Snake)
}

func JSONSchemaSnake2Camel(data []byte) []byte {
	return regex.JSONKey.ReplaceAllFunc(data, convertBytesSnake2Camel)
}

func convertBytesCamel2Snake(match []byte) []byte {
	return bytes.ToLower(regex.CamelCase.ReplaceAll(match, []byte(`${1}_${2}`)))
}

func convertBytesSnake2Camel(match []byte) []byte {
	str := string(match)
	key := str[1 : len(str)-2]
	resKey := toCamelCase(key)
	return []byte(`"` + resKey + `":`)
}

func toCamelCase(source string) string {
	strList := strings.Split(source, "_")
	for index, str := range strList {
		strList[index] = firstToUpper(str)
	}
	res := strings.Join(strList, "")
	return firstToLower(res)
}

func firstToUpper(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func firstToLower(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}
