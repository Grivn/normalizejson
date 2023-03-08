# GoJSON
Toolkits for golang to process JSON schema.

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
type FormatFunc func(item interface{}) (interface{}, error)
type FormatOption struct {
	FunctionName   string
	FormatFunction FormatFunc
}
```

The `FormatFunc` is used to convert JSON schema and the `FormatOption` has defined the name of each format_function. 

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
// DefaultOptions creates default options to format JSON schema.
var DefaultOptions = []FormatOption{
	{
		FunctionName:   FormatIntToString,
		FormatFunction: FormatDataIntToString,
	},
	{
		FunctionName:   FormatFloatToString,
		FormatFunction: FormatDataFloatToString,
	},
	{
		FunctionName:   FormatStringToInt,
		FormatFunction: FormatDataStringToInt,
	},
	{
		FunctionName:   FormatStringToFloat,
		FormatFunction: FormatDataStringToFloat,
	},
}

const (
	FormatIntToString = "int_to_string"
	FormatFloatToString  = "float_to_string"
	FormatStringToInt = "string_to_int"
	FormatStringToFloat  = "string_to_float"
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

In addition, you can create your own `FormatFunc` and add it into `FormatDataProvider` with `AddOptions` 
to implement your own method to normalize the JSON schema for designated part.

You can refer to [format_data_test.go](format_data_test.go) for details to create `FormatDataProvider` and normalize the JSON file.

## FormatKey

You can convert the key in JSON schema as below.

```go
type FormatFunc func(item interface{}) (interface{}, error)
type FormatKeyOption struct {
	From           string
	To             string
	FormatFunction FormatFunc
}
```

The `FormatKeyOption` has defined the method to convert JSON key, convert from `option.from` to `option.to` or take the `option.format_function` to convert the JSON key.

For each option, it can only have the `from` to `to` pair or have a `format_function`.
When `FormatKeyProvider` converts the JSON key, it tries to convert `from` to `to` at first, then invoke `format_function` to convert found JSON key.

You can refer to [format_key_test.go](format_key_test.go) for details to create `FormatKeyProvider` and normalize the JSON file.
