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
