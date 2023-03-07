package gojson

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatData(t *testing.T) {
	dir := "format_data"
	template, err := readTestData(dir, "config.json")
	if err != nil {
		panic(err)
	}

	source, err := readTestData(dir, "source.json")
	if err != nil {
		panic(err)
	}

	result, err := readTestData(dir, "result.json")
	if err != nil {
		panic(err)
	}

	options := append(DefaultOptions, createNilStringToBlankOption())
	formatted, err := JSONRawFormatData(source, template, options...)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(formatJSON(template)))
	fmt.Println(string(formatJSON(source)))
	fmt.Println(string(formatJSON(formatted)))
	assert.Equal(t, formatJSON(result), formatJSON(formatted))
}

func TestFormatDataProvider(t *testing.T) {
	dir := "format_data"
	template, err := readTestData(dir, "config.json")
	if err != nil {
		panic(err)
	}

	source, err := readTestData(dir, "source.json")
	if err != nil {
		panic(err)
	}

	result, err := readTestData(dir, "result.json")
	if err != nil {
		panic(err)
	}

	provider, err := NewDefaultFormatJSONProvider(template)
	if err != nil {
		panic(err)
	}

	provider.AddFormatOption(createNilStringToBlankOption())

	formatted, err := provider.FormatJSONData(source)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(formatJSON(template)))
	fmt.Println(string(formatJSON(source)))
	fmt.Println(string(formatJSON(formatted)))
	assert.Equal(t, formatJSON(result), formatJSON(formatted))
}

func createNilStringToBlankOption() JSONFormatOption {
	return JSONFormatOption{
		FunctionName:   "nil_string_to_blank",
		FormatFunction: nilStringToBlank,
	}
}

func nilStringToBlank(item interface{}, options ...interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}

	if str == "" {
		return "blank", nil
	}
	return item, nil
}
