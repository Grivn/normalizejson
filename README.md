# NormalizeJSON

<a href="https://godoc.org/github.com/Grivn/normalizejson"><img src="https://img.shields.io/badge/api-reference-pink.svg?style=flat-square" alt="GoDoc"></a>

NormalizeJSON is a Go package that provides a simple way to normalize the key/value in a JSON documents with a template.

This README is a quick overview of how to use NormalizeJSON. 

## Getting Started

### Installing

```shell
go get -u github.com/Grivn/normalizejson@latest
```

go version >= 1.18

### Create Provider

To take use of NormalizeJSON, you should create a `FormatProvider`.

```go
package main

import (
	"fmt"
	"github.com/Grivn/normalizejson"
)

func main() {
	provider, err := normalizejson.NewFormatProvider(nil)
	if err != nil {
		panic(fmt.Errorf("[EXAMPLE] create normalizejson provider failed: %s", err))
	}
	provider.Reset()
}
```

### Initiate Provider

Then, you should initiate the provider with `options` and `template`.

#### Options

The `options` will be used to normalize the key/value in JSON documents.

```go
package normalizejson

type FormatFuncType string

const (
	FormatFuncFormatData FormatFuncType = "format_function_type_format_data" // default type
	FormatFuncFormatKey                 = "format_function_type_format_key"
)

type FormatFunc func(item interface{}) (interface{}, error)

type FormatOption struct {
	FunctionType   FormatFuncType
	FunctionName   string
	FormatFunction FormatFunc
}
```

For each option, you should assign a `FormatFunc` in it.
This function will be used to normalize the key or value in JSON documents.

There are two types of options in NormalizeJSON:
- `normalizejson.FormatFuncFormatKey` (key-option)
- `normalizejson.FormatFuncFormatData` (data-option)

The key-options are used to normalize the JSON keys.
Each key should be normalized by every `FormatFunc` created from key-options.

The data-options are used to normalize the JSON values.
You should create a `template` to make statements about which function should be taken to normalize specific key's value. 

You can create key-options and value-options with methods `normalizejson.FormatKeyOption` and `normalizejson.FormatDataOption`.

Here's an example to initiate provider with options.

- To create key-options.

```go
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
```

- To create data-options.

```go
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
```

- To initiate provider. 

```go
package main

func main() {
	...

	options := append(FormatKeyOptions, FormatDataOptions...)
	provider.AddOptions(options...)
}
```

#### Template

To normalize the values in JSON document, you should create a template to state the function to use.

For example, to normalize [input.json](example/input.json) towards [output.json](example/output.json), you should create a template file [config.json](example/config.json). 

input.json

    {"data":{"description":1024,"id":"2","rate":"2.3","sub_data_list":[{"item1":12,"item2":"1.30","sub_data_list":[{"item1":"3","item2":"1.70","item3":"2","item4":1.2,"item5":"exist","type":"child"}],"type":"parent"},{"item1":"2","item2":"1.40","item3":"3","item4":"1.2000","item5":"","type":"child"}]}}

output.json

    {"data":{"description":"1024","id":"2","rate":2.3,"sub_data_list":[{"item1":12,"item2":"1.30","sub_data_list":[{"item1":3,"item2":"1.70","item3":"2","item4":1.2,"item5":"exist","type":"child"}],"type":"parent"},{"item1":2,"item2":"1.40","item3":"3","item4":1.2,"item5":"","type":"child"}]}}

config.json

```json
{
  "data": {
    "description": "to_string",
    "id": "__template.id",
    "rate": "to_float64",
    "sub_data_list": "__template.sub_data_list"
  },
  "id": "to_string",
  "sub_data": {
    "item1": "to_int64",
    "item3": "to_string",
    "item4": "to_float64",
    "sub_data_list": "__template.sub_data_list"
  },
  "sub_data_list": [
    "__template.sub_data"
  ]
}
```

This template file has defined 4 templates `data`, `id`, `sub_data`, and `sub_data_list`. 

The key-value in template show that the value in JSON document for current key should be processed by this format function.

    E.g.`{"data":{"description":"to_string"}}` 

    It means the value of `data.description` in JSON file should be processed by `FormatFunc` from data-options whose name is `to_string`.

There is a builtin style function name `__template.{{template_name}}`, which means we will process the value with template of `{{template_name}}`. 

    E.g. `{"sub_data":{"sub_data_list":"__template.sub_data_list"}}` 
    
    It means the value of `sub_data.sub_data_list` should be processed by template `sub_data_list`.

In addition, the statement like `["{{function_name}}"]` is used to describe the array structure in JSON document.
Each value in this array should be processed by function of `{{function_name}}`. 

    E.g. `{"sub_data_list":["__template.sub_data"]}` means the `sub_data_list` is an array of `sub_data`.
    And we should process each value of it with the `sub_data` template. 

To initiate the provider with `template`. 

```go
func main() {
	...
	
	if err = provider.UpdateTemplate(template); err != nil {
		panic(err)
	}
}
```

### Normalize JSON Schema

To normalize the JSON schema, just input the raw JSON document into `FormatJSONSchema`, then you can get the normalizedJSON. 

```go
func main() {
	...
	
	formattedJSON, err := provider.FormatJSONSchema(source)
	if err != nil {
		panic(err)
	}
}
```

## Example

You can take the [example](example) for details to use NormalizeJSON. 
It takes [template](example/config.json) to normalize the [input.json](example/input.json) to [output.json](example/output.json). 

```go
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
```

input.json

    {"data":{"description":1024,"id":"2","rate":"2.3","sub_data_list":[{"item1":12,"item2":"1.30","sub_data_list":[{"item1":"3","item2":"1.70","item3":"2","item4":1.2,"item5":"exist","type":"child"}],"type":"parent"},{"item1":"2","item2":"1.40","item3":"3","item4":"1.2000","item5":"","type":"child"}]}}

config.json

    {"data":{"description":"to_string","id":"__template.id","rate":"to_float64","sub_data_list":"__template.sub_data_list"},"id":"to_string","sub_data":{"item1":"to_int64","item3":"to_string","item4":"to_float64","sub_data_list":"__template.sub_data_list"},"sub_data_list":["__template.sub_data"]}

## TOOLs

You can refer to [Toolkits](TOOLKITs.md) to find some useful toolkits to process JSON schema. 
