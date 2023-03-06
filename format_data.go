package gojson

import (
	"encoding/json"
	"fmt"

	"github.com/andeya/ameda"
)

type FormatFunc func(map[string]interface{}, string) (map[string]interface{}, error)

type JSONFormatOption struct {
	FunctionName   string
	FormatFunction FormatFunc
}

func JSONRawFormatData(data []byte, conf interface{}, options ...JSONFormatOption) ([]byte, error) {
	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("unmarshal source data failed: %s", err)
	}

	formattedItem, err := formatItem(item, conf, createFormatFactory(options...))
	if err != nil {
		return nil, fmt.Errorf("format JSON data failed: %s", err)
	}
	return json.Marshal(formattedItem)
}

func DefaultJSONRawFormatData(data []byte, conf interface{}) ([]byte, error) {
	return JSONRawFormatData(data, conf, DefaultOptions...)
}

func createFormatFactory(options ...JSONFormatOption) map[string]FormatFunc {
	factory := make(map[string]FormatFunc)
	for _, option := range options {
		factory[option.FunctionName] = option.FormatFunction
	}
	return factory
}

// DefaultOptions creates default options to format JSON raw data.
var DefaultOptions = []JSONFormatOption{
	{
		FunctionName:   FormatNumberToString,
		FormatFunction: NumberToString,
	},
	{
		FunctionName:   FormatFloatToString,
		FormatFunction: FloatToString,
	},
	{
		FunctionName:   FormatStringToNumber,
		FormatFunction: StringToNumber,
	},
	{
		FunctionName:   FormatStringToFloat,
		FormatFunction: StringToFloat,
	},
}

const (
	FormatNumberToString = "number_to_string"
	FormatFloatToString  = "float_to_string"
	FormatStringToNumber = "string_to_number"
	FormatStringToFloat  = "string_to_float"
)

func NumberToString(jsonMap map[string]interface{}, key string) (map[string]interface{}, error) {
	rawID, ok := jsonMap[key]
	if !ok {
		// cannot find 'id' in panel's json map.
		return jsonMap, nil
	}

	float64ID, ok := rawID.(float64)
	if ok {
		int64ID, err := ameda.Float64ToInt64(float64ID)
		if err != nil {
			return jsonMap, err
		}
		jsonMap[key] = ameda.Int64ToString(int64ID)
		return jsonMap, nil
	}

	return jsonMap, nil
}

func FloatToString(jsonMap map[string]interface{}, key string) (map[string]interface{}, error) {
	rawID, ok := jsonMap[key]
	if !ok {
		// cannot find 'id' in panel's json map.
		return jsonMap, nil
	}

	float64ID, ok := rawID.(float64)
	if ok {
		jsonMap[key] = ameda.Float64ToString(float64ID)
		return jsonMap, nil
	}

	return jsonMap, nil
}

func StringToNumber(jsonMap map[string]interface{}, key string) (map[string]interface{}, error) {
	raw, ok := jsonMap[key]
	if !ok {
		return jsonMap, nil
	}
	str, ok := raw.(string)
	if !ok {
		return jsonMap, nil
	}
	if str == "" {
		jsonMap[key] = 0
		return jsonMap, nil
	}
	intValue, err := ameda.StringToInt64(str)
	if err != nil {
		return jsonMap, err
	}
	jsonMap[key] = intValue
	return jsonMap, nil
}

func StringToFloat(jsonMap map[string]interface{}, key string) (map[string]interface{}, error) {
	raw, ok := jsonMap[key]
	if !ok {
		return jsonMap, nil
	}
	str, ok := raw.(string)
	if !ok {
		return jsonMap, nil
	}
	if str == "" {
		jsonMap[key] = float64(0)
		return jsonMap, nil
	}
	floatValue, err := ameda.StringToFloat64(str)
	if err != nil {
		return jsonMap, err
	}
	jsonMap[key] = floatValue
	return jsonMap, nil
}
