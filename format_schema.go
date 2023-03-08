package gojson

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JSONSchemaFormat(data []byte, rawTemplate []byte, options ...FormatOption) ([]byte, error) {
	fsi, err := newFormatSchemaImpl(rawTemplate, options...)
	if err != nil {
		return nil, fmt.Errorf("format JSON data failed: %s", err)
	}
	return fsi.formatJSONSchema(data)
}

type formatSchemaImpl struct {
	formatKFunc map[string]FormatFunc
	formatVFunc map[string]FormatFunc
	templateMap map[string]interface{}
}

func newFormatSchemaImpl(rawTemplate []byte, options ...FormatOption) (*formatSchemaImpl, error) {
	fsi := &formatSchemaImpl{templateMap: make(map[string]interface{})}
	fsi.addOptions(options...)

	if len(rawTemplate) == 0 {
		return fsi, nil
	}

	if err := fsi.updateTemplate(rawTemplate); err != nil {
		return nil, err
	}
	return fsi, nil
}

func (fsi *formatSchemaImpl) reset() {
	fsi.formatKFunc = make(map[string]FormatFunc)
	fsi.formatVFunc = make(map[string]FormatFunc)
	fsi.templateMap = make(map[string]interface{})
}

func (fsi *formatSchemaImpl) updateTemplate(rawTemplate []byte) error {
	templateMap := make(map[string]interface{})
	if err := json.Unmarshal(rawTemplate, &templateMap); err != nil {
		return err
	}
	fsi.templateMap = templateMap
	return nil
}

func (fsi *formatSchemaImpl) addOptions(options ...FormatOption) {
	if fsi.formatKFunc == nil {
		fsi.formatKFunc = make(map[string]FormatFunc)
	}

	if fsi.formatVFunc == nil {
		fsi.formatVFunc = make(map[string]FormatFunc)
	}

	for _, option := range options {
		if option.FunctionType == FormatFuncFormatKey {
			fsi.formatKFunc[option.FunctionName] = option.FormatFunction
		} else {
			// format_function_type_format_key
			fsi.formatVFunc[option.FunctionName] = option.FormatFunction
		}
	}
}

func (fsi *formatSchemaImpl) formatJSONSchema(data []byte) ([]byte, error) {
	if len(fsi.formatKFunc) == 0 && len(fsi.formatVFunc) == 0 {
		return data, nil
	}

	var item interface{}
	if err := json.Unmarshal(data, &item); err != nil {
		return nil, fmt.Errorf("unmarshal source data failed: %s", err)
	}

	formattedItem, err := fsi.formatItem(item, nil)
	if err != nil {
		return nil, fmt.Errorf("format JSON data failed: %s", err)
	}
	return json.Marshal(formattedItem)
}

func (fsi *formatSchemaImpl) formatItem(item interface{}, template interface{}) (interface{}, error) {
	template = fsi.takeTemplate(template)

	switch v := item.(type) {
	case []interface{}:
		templateList, ok := template.([]interface{})
		if ok {
			return fsi.formatItemList(v, templateList)
		}
		return fsi.formatItemList(v, nil)
	case map[string]interface{}:
		templateMap, ok := template.(map[string]interface{})
		if ok {
			return fsi.formatItemMap(v, templateMap)
		}
		return fsi.formatItemMap(v, nil)
	default:
		expr, ok := template.(string)
		if !ok {
			return item, nil
		}

		f, ok := fsi.formatVFunc[expr]
		if ok {
			return f(item)
		}
		return item, nil
	}
}

func (fsi *formatSchemaImpl) formatItemList(itemList []interface{}, templateList []interface{}) ([]interface{}, error) {
	var template interface{}

	if len(templateList) > 0 {
		template = templateList[0]
	}

	for index, item := range itemList {
		formattedItem, err := fsi.formatItem(item, template)
		if err != nil {
			return itemList, err
		}
		itemList[index] = formattedItem
	}
	return itemList, nil
}

func (fsi *formatSchemaImpl) formatItemMap(itemMap map[string]interface{}, templateMap map[string]interface{}) (map[string]interface{}, error) {
	for key, item := range itemMap {
		// format JSON key at first.
		formattedKey, err := fsi.formatKey(key)
		if err != nil {
			return itemMap, err
		}

		// take the formatted key to find the template.
		var template interface{}
		if templateMap != nil {
			template = templateMap[formattedKey]
		} else {
			template = fsi.templateMap[formattedKey]
		}

		// format JSON value item with the selected template.
		formattedItem, err := fsi.formatItem(item, template)
		if err != nil {
			return itemMap, err
		}

		// update the item_map with the formatted JSON key and value.
		delete(itemMap, key)
		itemMap[formattedKey] = formattedItem
	}

	return itemMap, nil
}

func (fsi *formatSchemaImpl) formatKey(key string) (string, error) {
	for _, f := range fsi.formatKFunc {
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

func (fsi *formatSchemaImpl) takeTemplate(template interface{}) interface{} {
	expr, ok := template.(string)
	if !ok {
		return template
	}

	if !strings.HasPrefix(expr, formatDataTemplatePrefix) {
		return template
	}

	existTemplate, ok := fsi.templateMap[strings.TrimPrefix(expr, formatDataTemplatePrefix)]
	if !ok {
		return template
	}
	return existTemplate
}
