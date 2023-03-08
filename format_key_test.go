package gojson

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONSchemaFormatKey(t *testing.T) {
	dir := "format_key"

	source, err := readTestData(dir, "source.json")
	if err != nil {
		panic(err)
	}

	result, err := readTestData(dir, "result.json")
	if err != nil {
		panic(err)
	}

	formatted, err := JSONSchemaFormatKey(source, createFormatKeyTestOptions()...)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(result), formatJSON(formatted))
}

func TestFormatKeyProvider(t *testing.T) {
	dir := "format_key"

	source, err := readTestData(dir, "source.json")
	if err != nil {
		panic(err)
	}

	result, err := readTestData(dir, "result.json")
	if err != nil {
		panic(err)
	}

	formatter := NewFormatKeyProvider(createFormatKeyTestOptions()...)

	formatted, err := formatter.FormatJSONSchema(source)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(result), formatJSON(formatted))
}

func TestFunctionFormatKeyProvider(t *testing.T) {
	dir := "format_key"

	source, err := readTestData(dir, "source.json")
	if err != nil {
		panic(err)
	}

	result, err := readTestData(dir, "result.json")
	if err != nil {
		panic(err)
	}

	formatter := NewFormatKeyProvider(createFuncFormatKeyTestOptions()...)

	formatted, err := formatter.FormatJSONSchema(source)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(result), formatJSON(formatted))
}

func createFormatKeyTestOptions() []FormatKeyOption {
	return []FormatKeyOption{
		{
			From: "item1",
			To:   "result_item1",
		},
		{
			From: "item2",
			To:   "result_item2",
		},
		{
			From: "item3",
			To:   "result_item3",
		},
		{
			From: "item4",
			To:   "result_item4",
		},
		{
			From: "item5",
			To:   "result_item5",
		},
	}
}

func createFuncFormatKeyTestOptions() []FormatKeyOption {
	return []FormatKeyOption{
		{
			FormatFunction: formatKeyFunction,
		},
	}
}

func formatKeyFunction(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}

	if !strings.HasPrefix(str, "item") {
		return item, nil
	}

	return strings.ReplaceAll(str, "item", "result_item"), nil
}
