package gojson

import (
	"encoding/json"
)

func JSONSchemaFormatKey(data []byte, options ...FormatKeyOption) ([]byte, error) {
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

func (fki *formatKeyImpl) formatJSONSchema(data []byte) ([]byte, error) {
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

		formattedKey, ok := fki.formatKey(key)
		if ok {
			delete(itemMap, key)
		}

		itemMap[formattedKey] = formattedItem
	}

	return itemMap, nil
}

func (fki *formatKeyImpl) formatKey(key string) (string, bool) {
	to, ok := fki.formatKeyMap[key]
	if ok {
		return to, true
	}

	if fki.formatKeyFunc == nil {
		return key, false
	}

	raw, err := fki.formatKeyFunc(key)
	if err != nil {
		return key, false
	}

	to, ok = raw.(string)
	if ok {
		return to, true
	}

	return key, false
}
