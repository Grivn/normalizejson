package gojson

import (
	"github.com/andeya/ameda"
)

type FormatFunc func(item interface{}) (interface{}, error)

type FormatOption struct {
	FunctionName   string
	FormatFunction FormatFunc
}

const (
	FormatNumberToString = "number_to_string"
	FormatFloatToString  = "float_to_string"
	FormatStringToNumber = "string_to_number"
	FormatStringToFloat  = "string_to_float"
)

// DefaultFormatDataOptions creates default options to format JSON raw data.
var DefaultFormatDataOptions = []FormatOption{
	{
		FunctionName:   FormatNumberToString,
		FormatFunction: FormatDataNumberToString,
	},
	{
		FunctionName:   FormatFloatToString,
		FormatFunction: FormatDataFloatToString,
	},
	{
		FunctionName:   FormatStringToNumber,
		FormatFunction: FormatDataStringToNumber,
	},
	{
		FunctionName:   FormatStringToFloat,
		FormatFunction: FormatDataStringToFloat,
	},
}

func FormatDataNumberToString(item interface{}) (interface{}, error) {
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

func FormatDataStringToNumber(item interface{}) (interface{}, error) {
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
