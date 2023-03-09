package njson

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JSONSchemaFormatData(data []byte, rawTemplate []byte, options ...FormatOption) ([]byte, error) {
	fii, err := newFormatDataImpl(rawTemplate, options...)
	if err != nil {
		return nil, fmt.Errorf("format JSON data failed: %s", err)
	}
	return fii.formatJSONSchema(data)
}

func DefaultJSONSchemaFormatData(data []byte, rawTemplate []byte) ([]byte, error) {
	return JSONSchemaFormatData(data, rawTemplate, DefaultFormatDataOptions...)
}

type formatDataImpl struct {
	functionMap map[string]FormatFunc
	templateMap map[string]interface{}
}

const (
	formatDataTemplatePrefix = "__template."
)

func newFormatDataImpl(rawTemplate []byte, options ...FormatOption) (*formatDataImpl, error) {
	fii := &formatDataImpl{templateMap: make(map[string]interface{})}
	fii.addOptions(options...)

	if len(rawTemplate) == 0 {
		return fii, nil
	}

	if err := fii.updateTemplate(rawTemplate); err != nil {
		return nil, err
	}
	return fii, nil
}

func (fdi *formatDataImpl) reset() {
	fdi.functionMap = make(map[string]FormatFunc)
	fdi.templateMap = make(map[string]interface{})
}

func (fdi *formatDataImpl) updateTemplate(rawTemplate []byte) error {
	templateMap := make(map[string]interface{})
	if err := json.Unmarshal(rawTemplate, &templateMap); err != nil {
		return err
	}
	fdi.templateMap = templateMap
	return nil
}

func (fdi *formatDataImpl) addOptions(options ...FormatOption) {
	if fdi.functionMap == nil {
		fdi.functionMap = make(map[string]FormatFunc)
	}

	for _, option := range options {
		fdi.functionMap[option.FunctionName] = option.FormatFunction
	}
}

func (fdi *formatDataImpl) formatJSONSchema(data []byte) ([]byte, error) {
	if len(fdi.functionMap) == 0 {
		return data, nil
	}

	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("unmarshal source data failed: %s", err)
	}

	formattedItem, err := fdi.formatItem(item)
	if err != nil {
		return nil, fmt.Errorf("format JSON data failed: %s", err)
	}
	return json.Marshal(formattedItem)
}

func (fdi *formatDataImpl) formatItem(item interface{}) (interface{}, error) {
	switch v := item.(type) {
	case []interface{}:
		return fdi.formatItemList(v)
	case map[string]interface{}:
		return fdi.formatItemMap(v)
	default:
		return item, nil
	}
}

func (fdi *formatDataImpl) formatItemList(itemList []interface{}) ([]interface{}, error) {
	for index, item := range itemList {
		formattedItem, err := fdi.formatItem(item)
		if err != nil {
			return itemList, err
		}
		itemList[index] = formattedItem
	}
	return itemList, nil
}

func (fdi *formatDataImpl) formatItemMap(itemMap map[string]interface{}) (map[string]interface{}, error) {
	for key, item := range itemMap {
		template, ok := fdi.templateMap[key]
		if !ok {
			formattedItem, err := fdi.formatItem(item)
			if err != nil {
				return itemMap, err
			}
			itemMap[key] = formattedItem
		} else {
			formattedItem, err := fdi.formatItemByTemplate(item, template)
			if err != nil {
				return itemMap, err
			}
			itemMap[key] = formattedItem
		}
	}

	return itemMap, nil
}

func (fdi *formatDataImpl) formatItemByTemplate(item interface{}, template interface{}) (interface{}, error) {
	switch v := template.(type) {
	case string:
		if needTemplate(v) {
			return fdi.takeTemplate(item, v)
		}
		if f, ok := fdi.functionMap[v]; !ok {
			return item, nil
		} else {
			return f(item)
		}
	case []interface{}:
		itemList, ok := item.([]interface{})
		if !ok {
			return item, nil
		}
		return fdi.formatItemListByTemplate(itemList, v)
	case map[string]interface{}:
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return item, nil
		}
		return fdi.formatItemMapByTemplate(itemMap, v)
	default:
		return item, nil
	}
}

func (fdi *formatDataImpl) formatItemListByTemplate(itemList []interface{}, templateList []interface{}) ([]interface{}, error) {
	if len(templateList) == 0 {
		return itemList, nil
	}
	for index, item := range itemList {
		formattedItem, err := fdi.formatItemByTemplate(item, templateList[0])
		if err != nil {
			return itemList, err
		}
		itemList[index] = formattedItem
	}
	return itemList, nil
}

func (fdi *formatDataImpl) formatItemMapByTemplate(itemMap map[string]interface{}, templateMap map[string]interface{}) (map[string]interface{}, error) {
	for key, template := range templateMap {
		item, exist := itemMap[key]
		if !exist {
			continue
		}

		formattedItem, err := fdi.formatItemByTemplate(item, template)
		if err != nil {
			return itemMap, err
		}
		itemMap[key] = formattedItem
	}
	return itemMap, nil
}

func (fdi *formatDataImpl) takeTemplate(item interface{}, expr string) (interface{}, error) {
	templateKey := strings.TrimPrefix(expr, formatDataTemplatePrefix)
	template, ok := fdi.templateMap[templateKey]
	if !ok {
		return item, nil
	}
	return fdi.formatItemByTemplate(item, template)
}

func needTemplate(function string) bool {
	return strings.HasPrefix(function, formatDataTemplatePrefix)
}
