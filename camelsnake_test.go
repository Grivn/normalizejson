package gojson

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONRawCamel2Snake(t *testing.T) {
	dir := "camelsnake"

	input, err := readTestData(dir, "camel.json")
	if err != nil {
		panic(err)
	}

	output, err := readTestData(dir, "snake.json")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, output, JSONRawCamel2Snake(input))
}

func TestJSONRawSnake2Camel(t *testing.T) {
	dir := "camelsnake"

	input, err := readTestData(dir, "snake.json")
	if err != nil {
		panic(err)
	}

	output, err := readTestData(dir, "camel.json")
	if err != nil {
		panic(err)
	}

	assert.Equal(t, output, JSONRawSnake2Camel(input))
}
