package njson

import (
	"encoding/json"
	"fmt"
)

func JSONSchemaFormatKey(data []byte, options ...FormatOption) ([]byte, error) {
	formatter := newFormatKeyImpl(options...)
	return formatter.formatJSONSchema(data)
}

type formatKeyImpl struct {
	functionMap map[string]FormatFunc
}

func newFormatKeyImpl(options ...FormatOption) *formatKeyImpl {
	fki := &formatKeyImpl{}
	fki.addOptions(options...)
	return fki
}

func (fki *formatKeyImpl) formatJSONSchema(data []byte) ([]byte, error) {
	if len(fki.functionMap) == 0 {
		return data, nil
	}

	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return data, err
	}

	formattedItem, err := fki.formatItem(item)
	if err != nil {
		return data, err
	}

	return json.Marshal(formattedItem)
}

func (fki *formatKeyImpl) reset() {
	fki.functionMap = make(map[string]FormatFunc)
}

func (fki *formatKeyImpl) addOptions(options ...FormatOption) {
	if fki.functionMap == nil {
		fki.functionMap = make(map[string]FormatFunc)
	}

	for _, option := range options {
		fki.functionMap[option.FunctionName] = option.FormatFunction
	}
}

func (fki *formatKeyImpl) formatItem(item interface{}) (interface{}, error) {
	switch v := item.(type) {
	case []interface{}:
		return fki.formatItemList(v)
	case map[string]interface{}:
		return fki.formatItemMap(v)
	default:
		return item, nil
	}
}

func (fki *formatKeyImpl) formatItemList(itemList []interface{}) ([]interface{}, error) {
	for index, item := range itemList {
		formattedItem, err := fki.formatItem(item)
		if err != nil {
			return itemList, err
		}
		itemList[index] = formattedItem
	}
	return itemList, nil
}

func (fki *formatKeyImpl) formatItemMap(itemMap map[string]interface{}) (map[string]interface{}, error) {
	for key, item := range itemMap {
		formattedItem, err := fki.formatItem(item)
		if err != nil {
			return itemMap, err
		}

		formattedKey, err := fki.formatKey(key)
		if err != nil {
			return itemMap, err
		}

		delete(itemMap, key)
		itemMap[formattedKey] = formattedItem
	}

	return itemMap, nil
}

func (fki *formatKeyImpl) formatKey(key string) (string, error) {
	for _, f := range fki.functionMap {
		formatted, err := f(key)
		if err != nil {
			return key, err
		}

		formattedKey, ok := formatted.(string)
		if !ok {
			return key, fmt.Errorf("illegal converted type")
		}
		key = formattedKey
	}
	return key, nil
}
