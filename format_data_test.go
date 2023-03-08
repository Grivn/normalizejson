package gojson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Schema struct {
	Data Data `json:"data"`
}

type Data struct {
	ID          string      `json:"id"`
	Description string      `json:"description"`
	Rate        float64     `json:"rate"`
	SubDataList SubDataList `json:"sub_data_list"`
}

type SubDataList []SubData
type SubData struct {
	Type        string      `json:"type"`
	Item1       int         `json:"item1"`
	Item2       string      `json:"item2"`
	Item3       string      `json:"item3,omitempty"`
	Item4       float64     `json:"item4,omitempty"`
	Item5       string      `json:"item5,omitempty"`
	SubDataList SubDataList `json:"sub_data_list,omitempty"`
}

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

func TestFormatDataProviderAndUpdateTemplate(t *testing.T) {
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
