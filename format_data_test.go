package gojson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatData(t *testing.T) {
	dir := "format_data"
	confRaw, err := readTestData(dir, "config.json")
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

	cleanResult, err := removeJSONBlankAndBreak(result)
	if err != nil {
		panic(err)
	}

	var conf interface{}
	if err = json.Unmarshal(confRaw, &conf); err != nil {
		panic(err)
	}

	formatted, err := DefaultJSONRawFormatData(source, conf)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, cleanResult, formatted)
}
