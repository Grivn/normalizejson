package gojson

import (
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
)

// DefaultFormatDataOptions creates default options to format JSON raw data.
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
