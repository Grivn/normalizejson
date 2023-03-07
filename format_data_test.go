package gojson

import (
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

	formatted, err := DefaultJSONRawFormatData(source, template)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, formatJSON(result), formatJSON(formatted))
}
