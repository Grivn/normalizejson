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

	formatted, err := JSONSchemaFormatKey(source, createFuncFormatKeyOption())
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

	provider := NewFormatKeyProvider()
	provider.AddOptions(createFuncFormatKeyOption())

	formatted, err := provider.FormatJSONSchema(source)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(result), formatJSON(formatted))
}

func TestFormatKeyProviderCamelToSnake(t *testing.T) {
	dir := "camel_snake"

	camel, err := readTestData(dir, "camel.json")
	if err != nil {
		panic(err)
	}

	snake, err := readTestData(dir, "snake.json")
	if err != nil {
		panic(err)
	}

	provider := NewFormatKeyProvider()
	provider.AddOptions(FormatKeyOption(FormatCamelToSnake, FormatKeyCamelToSnake))

	formattedSnake, err := provider.FormatJSONSchema(camel)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(snake), formatJSON(formattedSnake))
}

func TestFormatKeyProviderSnakeToCamel(t *testing.T) {
	dir := "camel_snake"

	camel, err := readTestData(dir, "camel.json")
	if err != nil {
		panic(err)
	}

	snake, err := readTestData(dir, "snake.json")
	if err != nil {
		panic(err)
	}

	provider := NewFormatKeyProvider()
	provider.AddOptions(FormatKeyOption(FormatSnakeToCamel, FormatKeySnakeToCamel))

	formattedCamel, err := provider.FormatJSONSchema(snake)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(camel), formatJSON(formattedCamel))
}

const (
	formatAppendSuffix = "append_suffix"
)

func createFuncFormatKeyOption() FormatOption {
	return FormatKeyOption(formatAppendSuffix, formatKeyAppendSuffix)
}

func formatKeyAppendSuffix(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}

	if !strings.HasPrefix(str, "item") {
		return item, nil
	}

	return strings.ReplaceAll(str, "item", "result_item"), nil
}
