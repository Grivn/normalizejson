package njson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatData(t *testing.T) {
	s := Schema{}
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

	formatted, err := DefaultJSONSchemaFormatData(source, template)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(result), formatJSON(formatted))

	assert.NotNil(t, json.Unmarshal(source, &s))
	assert.Nil(t, json.Unmarshal(formatted, &s))
}

func TestFormatDataProviderAndUpdateTemplateAndReset(t *testing.T) {
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

	provider, err := NewDefaultFormatDataProvider(nil)
	if err != nil {
		panic(err)
	}

	if err = provider.UpdateTemplate(template); err != nil {
		panic(err)
	}

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

func createNilStringToBlankOption() FormatOption {
	return FormatDataOption("nil_string_to_blank", formatDataNilStringToBlank)
}

func formatDataNilStringToBlank(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}

	if str == "" {
		return "blank", nil
	}
	return item, nil
}
