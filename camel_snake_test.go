package gojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONSchemaCamel2Snake(t *testing.T) {
	dir := "camel_snake"

	input, err := readTestData(dir, "camel.json")
	if err != nil {
		panic(err)
	}

	output, err := readTestData(dir, "snake.json")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, output, JSONSchemaCamel2Snake(input))
}

func TestJSONSchemaSnake2Camel(t *testing.T) {
	dir := "camel_snake"

	input, err := readTestData(dir, "snake.json")
	if err != nil {
		panic(err)
	}

	output, err := readTestData(dir, "camel.json")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, output, JSONSchemaSnake2Camel(input))
}
