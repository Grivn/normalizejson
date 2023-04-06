package normalizejson

import (
	"strings"

	"github.com/Grivn/normalizejson/regex"
	"github.com/spf13/cast"
)

type FormatFuncType string

const (
	FormatFuncFormatData FormatFuncType = "format_function_type_format_data" // default type
	FormatFuncFormatKey                 = "format_function_type_format_key"
)

type FormatFunc func(item interface{}) (interface{}, error)

type FormatOption struct {
	FunctionType   FormatFuncType
	FunctionName   string
	FormatFunction FormatFunc
	RetainKey      bool
}

func FormatDataOption(funcName string, formatFunc FormatFunc) FormatOption {
	return FormatOption{
		FunctionType:   FormatFuncFormatData,
		FunctionName:   funcName,
		FormatFunction: formatFunc,
	}
}

func FormatKeyOption(funcName string, formatFunc FormatFunc) FormatOption {
	return FormatOption{
		FunctionType:   FormatFuncFormatKey,
		FunctionName:   funcName,
		FormatFunction: formatFunc,
	}
}

const (
	FormatToInt64   = "to_int64"
	FormatToFloat64 = "to_float64"
	FormatToString  = "to_string"
	FormatToBool    = "to_bool"

	FormatCamelToSnake = "camel_to_snake"
	FormatSnakeToCamel = "snake_to_camel"
)

var DefaultFormatDataOptions = []FormatOption{
	FormatDataOption(FormatToInt64, FormatDataToInt64),
	FormatDataOption(FormatToFloat64, FormatDataToFloat64),
	FormatDataOption(FormatToString, FormatDataToString),
	FormatDataOption(FormatToBool, FormatDataToBool),
}

func FormatDataToString(item interface{}) (interface{}, error) {
	return cast.ToStringE(item)
}

func FormatDataToInt64(item interface{}) (interface{}, error) {
	return cast.ToInt64E(item)
}

func FormatDataToFloat64(item interface{}) (interface{}, error) {
	return cast.ToFloat64E(item)
}

func FormatDataToBool(item interface{}) (interface{}, error) {
	return cast.ToBoolE(item)
}

func FormatKeyCamelToSnake(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}
	return strings.ToLower(regex.CamelCase.ReplaceAllString(str, `${1}_${2}`)), nil
}

func FormatKeySnakeToCamel(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}

	strList := strings.Split(str, "_")
	for index, str := range strList {
		strList[index] = firstToUpper(str)
	}
	res := strings.Join(strList, "")
	return firstToLower(res), nil
}
