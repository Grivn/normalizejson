# GoJSON
Normalize the key and value in JSON schema to specific type.

## Prepare

```shell
go get github.com/Grivn/gojson@latest
```

## Camel-Case & Snake-Case

You can convert the JSON schema between camel-case and snake-case as below.

- `JSONSchemaCamel2Snake` converts JSON schema from camel-case to snake-case.
- `JSONSchemaSnake2Camel` converts JSON schema from snake-case to camel-case.

## FormatData

You can convert the JSON schema as below.

```go
package main
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

The `FormatFunc` is used to convert JSON schema and the `FormatOption` has defined the name of each format_function. 
To convert the data values of JSON schema, we should define the `FormatFuncType` as `FormatFuncFormatData`. 

For instance, there is a source JSON file as below.

```json
{
  "data": {
    "description": 1024,
    "id": "2",
    "rate": "2.3",
    "sub_data_list": [
      {
        "item1": 12,
        "item2": "1.30",
        "sub_data_list": [
          {
            "item1": "3",
            "item2": "1.70",
            "item3": "2",
            "item4": 1.2,
            "item5": "exist",
            "type": "child"
          }
        ],
        "type": "parent"
      },
      {
        "item1": "2",
        "item2": "1.40",
        "item3": "3",
        "item4": "1.2000",
        "item5": "",
        "type": "child"
      }
    ]
  }
}
```

Then we try to normalize the source JSON file to the standard schema.

```go
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
```

Here, we need to convert that
`data.id` to string,
`sub_data.item1` to number,
`sub_data.item3` to string,
`sub_data.item4` to float,
`sub_data.item5` to another format of string.
And the components in JSON might be nested, such as `sub_data_list` in `sub_data`.

So that, we need to create template and create `FormatDataProvider` with options for `FormatFunc`.

The example template is as following.
The values `int_to_string/float_to_string/string_to_int/string_to_float` refer to the names of `format_function` created by options.
The values `__template.id/__template.sub_data_list/__template.sub_data` refer to the structure of current JSON schema. 

```json
{
  "data": {
    "description": "int_to_string",
    "id": "__template.id",
    "rate": "string_to_float",
    "sub_data_list": "__template.sub_data_list"
  },
  "id": "int_to_string",
  "sub_data": {
    "item1": "string_to_int",
    "item3": "float_to_string",
    "item4": "string_to_float",
    "item5": "nil_string_to_blank",
    "sub_data_list": "__template.sub_data_list"
  },
  "sub_data_list": [
    "__template.sub_data"
  ]
}
```

The default `FormatFunc` options are as follows.

```go
func FormatDataOption(funcName string, formatFunc FormatFunc) FormatOption {
	return FormatOption{
		FunctionType:   FormatFuncFormatData,
		FunctionName:   funcName,
		FormatFunction: formatFunc,
	}
}

const (
	FormatIntToString   = "int_to_string"
	FormatFloatToString = "float_to_string"
	FormatStringToInt   = "string_to_int"
	FormatStringToFloat = "string_to_float"
)

func FormatDataIntToString(item interface{}) (interface{}, error) {
	float64ID, ok := item.(float64)
	if ok {
		int64ID, err := ameda.Float64ToInt64(float64ID)
		if err != nil {
			return item, err
		}
		return ameda.Int64ToString(int64ID), nil
	}
	return item, nil
}

func FormatDataFloatToString(item interface{}) (interface{}, error) {
	float64ID, ok := item.(float64)
	if ok {
		return ameda.Float64ToString(float64ID), nil
	}
	return item, nil
}

func FormatDataStringToInt(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}
	if str == "" {
		return 0, nil
	}
	intValue, err := ameda.StringToInt64(str)
	if err != nil {
		return item, err
	}
	return intValue, nil
}

func FormatDataStringToFloat(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}
	if str == "" {
		return float64(0), nil
	}
	floatValue, err := ameda.StringToFloat64(str)
	if err != nil {
		return item, err
	}
	return floatValue, nil
}
```

```go
var DefaultFormatDataOptions = []FormatOption{
	FormatDataOption(FormatIntToString, FormatDataIntToString),
	FormatDataOption(FormatFloatToString, FormatDataFloatToString),
	FormatDataOption(FormatStringToInt, FormatDataStringToInt),
	FormatDataOption(FormatStringToFloat, FormatDataStringToFloat),
}
```

In addition, you can create your own `FormatFunc` and add it into `FormatDataProvider` with `AddOptions` 
to implement your own method to normalize the JSON schema for designated part.

You can refer to [format_data_test.go](format_data_test.go) for details to create `FormatProvider` and normalize the value data in JSON file.

## FormatKey

You can convert the key in JSON schema as below.

```go
package main
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

The `FormatOption` has defined the method to convert JSON key.
To convert JSON key, we need to define the `FormatFuncType` as `FormatFuncFormatKey`. 

```go
func FormatKeyOption(funcName string, formatFunc FormatFunc) FormatOption {
	return FormatOption{
		FunctionType:   FormatFuncFormatKey,
		FunctionName:   funcName,
		FormatFunction: formatFunc,
	}
}

const (
	FormatCamelToSnake = "camel_to_snake"
)

func FormatKeyCamelToSnake(item interface{}) (interface{}, error) {
	str, ok := item.(string)
	if !ok {
		return item, nil
	}
	return strings.ToLower(regex.CamelCase.ReplaceAllString(str, `${1}_${2}`)), nil
}
```

```go
var option = FormatKeyOption(FormatSnakeToCamel, FormatKeySnakeToCamel)
```

When `FormatKeyProvider` converts the JSON key, it tries to invoke each `FormatFunc` from options to convert found JSON key.

You can refer to [format_key_test.go](format_key_test.go) for details to create `FormatProvider` and normalize the key in JSON file.

## FormatSchema

To convert JSON key and value at the same time, you can create `FormatSchemaProvider`. 

You should define every option you need and create the template to convert JSON value data. 

You can refer to [format_schema_test.go](format_schema_test.go) for details to create `FormatProvider` and normalize the key and value data JSON file.

