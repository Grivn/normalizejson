package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Grivn/normalizejson"
	"github.com/spf13/cast"
)

func main() {
	provider, err := normalizejson.NewFormatProvider(nil)
	if err != nil {
		panic(fmt.Errorf("[EXAMPLE] create normalizejson provider failed: %s", err))
	}

	provider.Reset()

	options := append(FormatKeyOptions, FormatDataOptions...)
	provider.AddOptions(options...)

	source, err := os.ReadFile("input.json")
	if err != nil {
		panic(err)
	}

	template, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err = provider.UpdateTemplate(template); err != nil {
		panic(err)
	}

	formattedJSON, err := provider.FormatJSONSchema(source)
	if err != nil {
		panic(err)
	}

	if _, err = os.Create("output.json"); err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err = json.Indent(&buf, formattedJSON, "", "  "); err != nil {
		panic(err)
	}

	if err = os.WriteFile("output.json", buf.Bytes(), 0777); err != nil {
		panic(err)
	}
}

const (
	FormatCamelToSnake = "camel_to_snake"
)

var FormatKeyOptions = []normalizejson.FormatOption{
	normalizejson.FormatKeyOption(FormatCamelToSnake, FormatKeyCamelToSnake),
}

var regexCamelCaseJSONKey = regexp.MustCompile(`\"(\w+)\":`)

func FormatKeyCamelToSnake(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}
	return strings.ToLower(regexCamelCaseJSONKey.ReplaceAllString(str, `${1}_${2}`)), nil
}

var FormatDataOptions = []normalizejson.FormatOption{
	normalizejson.FormatDataOption(FormatToInt64, FormatDataToInt64),
	normalizejson.FormatDataOption(FormatToFloat64, FormatDataToFloat64),
	normalizejson.FormatDataOption(FormatToString, FormatDataToString),
	normalizejson.FormatDataOption(FormatToBool, FormatDataToBool),
}

const (
	FormatToInt64   = "to_int64"
	FormatToFloat64 = "to_float64"
	FormatToString  = "to_string"
	FormatToBool    = "to_bool"
)

func FormatDataToString(item interface{}) (interface{}, error) {
	return cast.ToStringE(item)
}

func FormatDataToInt64(item interface{}) (interface{}, error) {
	return cast.ToInt64E(item)
}

func FormatDataToFloat64(item interface{}) (interface{}, error) {
	return cast.ToFloat64E(item)
}

func FormatDataToBool(item interface{}) (interface{}, error) {
	return cast.ToBoolE(item)
}

func printJSON(raw []byte) {
	fmt.Println(string(formatJSON(raw)))
}

func formatJSON(raw []byte) []byte {
	formattedJSON, _ := removeJSONBlankAndBreak(raw)
	return formattedJSON
}

func removeJSONBlankAndBreak(raw []byte) ([]byte, error) {
	var item interface{}
	if err := json.Unmarshal(raw, &item); err != nil {
		return nil, err
	}
	return json.Marshal(item)
}
