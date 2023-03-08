package gojson

import (
	"strings"

	"github.com/Grivn/gojson/regex"
	"github.com/andeya/ameda"
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
	FormatIntToString   = "int_to_string"
	FormatFloatToString = "float_to_string"
	FormatStringToInt   = "string_to_int"
	FormatStringToFloat = "string_to_float"

	FormatCamelToSnake = "camel_to_snake"
	FormatSnakeToCamel = "snake_to_camel"
)

var DefaultFormatDataOptions = []FormatOption{
	FormatDataOption(FormatIntToString, FormatDataIntToString),
	FormatDataOption(FormatFloatToString, FormatDataFloatToString),
	FormatDataOption(FormatStringToInt, FormatDataStringToInt),
	FormatDataOption(FormatStringToFloat, FormatDataStringToFloat),
}

func FormatDataIntToString(item interface{}) (interface{}, error) {
	float64ID, ok := item.(float64)
	if ok {
		int64ID, err := ameda.Float64ToInt64(float64ID)
		if err != nil {
			return item, err
		}
		return ameda.Int64ToString(int64ID), nil
	}
	return item, nil
}

func FormatDataFloatToString(item interface{}) (interface{}, error) {
	float64ID, ok := item.(float64)
	if ok {
		return ameda.Float64ToString(float64ID), nil
	}
	return item, nil
}

func FormatDataStringToInt(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}
	if str == "" {
		return 0, nil
	}
	intValue, err := ameda.StringToInt64(str)
	if err != nil {
		return item, err
	}
	return intValue, nil
}

func FormatDataStringToFloat(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}
	if str == "" {
		return float64(0), nil
	}
	floatValue, err := ameda.StringToFloat64(str)
	if err != nil {
		return item, err
	}
	return floatValue, nil
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
