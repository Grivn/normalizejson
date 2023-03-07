package gojson

import (
	"github.com/Grivn/gojson/regex"
)

func JSONSchemaFormatKey(data []byte, options ...FormatKeyOption) []byte {
	formatter := newFormatKeyImpl(options...)
	return formatter.formatJSONSchema(data)
}

type FormatKeyOption struct {
	From           string
	To             string
	FormatFunction FormatFunc
}

type formatKeyImpl struct {
	formatKeyMap  map[string]string
	formatKeyFunc FormatFunc
}

func newFormatKeyImpl(options ...FormatKeyOption) *formatKeyImpl {
	fki := &formatKeyImpl{}
	fki.addOptions(options...)
	return fki
}

func (fki *formatKeyImpl) formatJSONSchema(data []byte) []byte {
	return regex.JSONKey.ReplaceAllFunc(data, fki.formatKey)
}

func (fki *formatKeyImpl) formatKey(match []byte) []byte {
	key := readKey(match)
	to, ok := fki.formatKeyMap[key]
	if ok {
		return createKey(to)
	}

	if fki.formatKeyFunc == nil {
		return match
	}

	raw, err := fki.formatKeyFunc(key)
	if err != nil {
		return match
	}

	to, ok = raw.(string)
	if ok {
		return createKey(to)
	}

	return match
}

func (fki *formatKeyImpl) addOptions(options ...FormatKeyOption) {
	if fki.formatKeyMap == nil {
		fki.formatKeyMap = make(map[string]string)
	}

	for _, option := range options {
		fki.formatKeyMap[option.From] = option.To

		if option.FormatFunction == nil {
			continue
		}
		fki.formatKeyFunc = option.FormatFunction
	}
}

func readKey(raw []byte) string {
	str := string(raw)
	key := str[1 : len(str)-2]
	return key
}

func createKey(key string) []byte {
	return []byte(`"` + key + `":`)
}
