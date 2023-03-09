package njson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSchema(t *testing.T) {
	s := Schema{}
	dir := "format_schema"
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

	options := append(DefaultFormatDataOptions, FormatKeyOption(FormatCamelToSnake, FormatKeyCamelToSnake))
	formatted, err := JSONSchemaFormat(source, template, options...)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(result), formatJSON(formatted))

	assert.NotNil(t, json.Unmarshal(source, &s))
	assert.Nil(t, json.Unmarshal(formatted, &s))
}

func TestFormatSchemaProviderAndUpdateTemplateAndReset(t *testing.T) {
	dir := "format_schema"
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

	provider, err := NewDefaultFormatSchemaProvider(nil)
	if err != nil {
		panic(err)
	}

	if err = provider.UpdateTemplate(template); err != nil {
		panic(err)
	}

	// add option to format JSON key.
	provider.AddOptions(FormatKeyOption(FormatCamelToSnake, FormatKeyCamelToSnake))

	formatted, err := provider.FormatJSONSchema(source)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(result), formatJSON(formatted))

	templateToBlank, err := readTestData(dir, "config_to_blank.json")
	if err != nil {
		panic(err)
	}

	resultToBlank, err := readTestData(dir, "result_to_blank.json")
	if err != nil {
		panic(err)
	}

	provider.AddOptions(createNilStringToBlankOption())

	if err = provider.UpdateTemplate(templateToBlank); err != nil {
		panic(err)
	}

	formattedToBlank, err := provider.FormatJSONSchema(source)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(resultToBlank), formatJSON(formattedToBlank))

	provider.Reset()
	if err = provider.UpdateTemplate(templateToBlank); err != nil {
		panic(err)
	}

	formattedFailed, err := provider.FormatJSONSchema(source)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(source), formatJSON(formattedFailed))
}
