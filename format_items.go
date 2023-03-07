package gojson

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	formatDataTemplatePrefix = "__template."
)

type formatItemsImpl struct {
	functionMap map[string]FormatFunc
	templateMap map[string]interface{}
}

func newFormatItemImpl(rawTemplate []byte, options ...JSONFormatOption) (*formatItemsImpl, error) {
	fii := &formatItemsImpl{
		functionMap: createFormatFactory(options...),
		templateMap: make(map[string]interface{}),
	}

	if len(rawTemplate) == 0 {
		return fii, nil
	}

	if err := fii.updateTemplate(rawTemplate); err != nil {
		return nil, err
	}
	return fii, nil
}

func (fii *formatItemsImpl) updateTemplate(rawTemplate []byte) error {
	templateMap := make(map[string]interface{})
	if err := json.Unmarshal(rawTemplate, &templateMap); err != nil {
		return err
	}
	fii.templateMap = templateMap
	return nil
}

func (fii *formatItemsImpl) addFormatOption(options ...JSONFormatOption) {
	if fii.functionMap == nil {
		fii.functionMap = make(map[string]FormatFunc)
	}

	for _, option := range options {
		fii.functionMap[option.FunctionName] = option.FormatFunction
	}
}

func (fii *formatItemsImpl) formatJSONData(data []byte) ([]byte, error) {
	if len(fii.functionMap) == 0 {
		return data, nil
	}

	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("unmarshal source data failed: %s", err)
	}

	formattedItem, err := fii.formatItem(item)
	if err != nil {
		return nil, fmt.Errorf("format JSON data failed: %s", err)
	}
	return json.Marshal(formattedItem)
}

func (fii *formatItemsImpl) formatItem(item interface{}) (interface{}, error) {
	switch v := item.(type) {
	case []interface{}:
		return fii.formatItemList(v)
	case map[string]interface{}:
		return fii.formatItemMap(v)
	default:
		return item, nil
	}
}

func (fii *formatItemsImpl) formatItemList(itemList []interface{}) ([]interface{}, error) {
	for index, item := range itemList {
		formattedItem, err := fii.formatItem(item)
		if err != nil {
			return itemList, err
		}
		itemList[index] = formattedItem
	}
	return itemList, nil
}

func (fii *formatItemsImpl) formatItemMap(itemMap map[string]interface{}) (map[string]interface{}, error) {
	for key, item := range itemMap {
		template, ok := fii.templateMap[key]
		if !ok {
			formattedItem, err := fii.formatItem(item)
			if err != nil {
				return itemMap, err
			}
			itemMap[key] = formattedItem
		} else {
			formattedItem, err := fii.formatItemByTemplate(item, template)
			if err != nil {
				return itemMap, err
			}
			itemMap[key] = formattedItem
		}
	}

	return itemMap, nil
}

func (fii *formatItemsImpl) formatItemByTemplate(item interface{}, template interface{}) (interface{}, error) {
	switch v := template.(type) {
	case string:
		if needTemplate(v) {
			return fii.takeTemplate(item, v)
		}
		if f, ok := fii.functionMap[v]; !ok {
			return item, nil
		} else {
			return f(item)
		}
	case []interface{}:
		itemList, ok := item.([]interface{})
		if !ok {
			return item, nil
		}
		return fii.formatItemListByTemplate(itemList, v)
	case map[string]interface{}:
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return item, nil
		}
		return fii.formatItemMapByTemplate(itemMap, v)
	default:
		return item, nil
	}
}

func (fii *formatItemsImpl) formatItemListByTemplate(itemList []interface{}, templateList []interface{}) ([]interface{}, error) {
	if len(templateList) == 0 {
		return itemList, nil
	}
	for index, item := range itemList {
		formattedItem, err := fii.formatItemByTemplate(item, templateList[0])
		if err != nil {
			return itemList, err
		}
		itemList[index] = formattedItem
	}
	return itemList, nil
}

func (fii *formatItemsImpl) formatItemMapByTemplate(itemMap map[string]interface{}, templateMap map[string]interface{}) (map[string]interface{}, error) {
	for key, template := range templateMap {
		item, exist := itemMap[key]
		if !exist {
			continue
		}

		formattedItem, err := fii.formatItemByTemplate(item, template)
		if err != nil {
			return itemMap, err
		}
		itemMap[key] = formattedItem
	}
	return itemMap, nil
}

func (fii *formatItemsImpl) takeTemplate(item interface{}, expr string) (interface{}, error) {
	templateKey := strings.TrimPrefix(expr, formatDataTemplatePrefix)
	template, ok := fii.templateMap[templateKey]
	if !ok {
		return item, nil
	}
	return fii.formatItemByTemplate(item, template)
}

func needTemplate(function string) bool {
	return strings.HasPrefix(function, formatDataTemplatePrefix)
}
