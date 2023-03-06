package gojson

import "fmt"

func formatItem(item interface{}, conf interface{}, funcMap map[string]FormatFunc) (interface{}, error) {
	switch v := conf.(type) {
	case []interface{}:
		itemList, ok := item.([]interface{})
		if !ok {
			return item, nil
		}
		return formatItemList(itemList, v, funcMap)
	case map[string]interface{}:
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return item, nil
		}
		return formatItemMap(itemMap, v, funcMap)
	default:
		return item, nil
	}
}

func formatItemList(itemList []interface{}, confList []interface{}, funcMap map[string]FormatFunc) ([]interface{}, error) {
	if len(confList) == 0 {
		return itemList, nil
	}
	for index, item := range itemList {
		formatdItem, err := formatItem(item, confList[0], funcMap)
		if err != nil {
			return itemList, err
		}
		itemList[index] = formatdItem
	}
	return itemList, nil
}

func formatItemMap(itemMap map[string]interface{}, confMap map[string]interface{}, funcMap map[string]FormatFunc) (map[string]interface{}, error) {
	for key, conf := range confMap {
		item, exist := itemMap[key]
		if !exist {
			continue
		}

		switch v := conf.(type) {
		case string:
			f, ok := funcMap[v]
			if !ok {
				return itemMap, fmt.Errorf("unknown format function: %s", v)
			}

			formattedMap, err := f(itemMap, key)
			if err != nil {
				return itemMap, err
			}
			itemMap = formattedMap
		case map[string]interface{}:
			if subMap, ok := item.(map[string]interface{}); ok {
				formattedMap, err := formatItemMap(subMap, v, funcMap)
				if err != nil {
					return itemMap, err
				}
				itemMap[key] = formattedMap
			}
		case []interface{}:
			if subList, ok := item.([]interface{}); ok {
				formattedMap, err := formatItemList(subList, v, funcMap)
				if err != nil {
					return itemMap, err
				}
				itemMap[key] = formattedMap
			}
		}
	}
	return itemMap, nil
}
