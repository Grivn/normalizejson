package gojson

import (
	"encoding/json"
	"os"
)

func readTestData(dir, file string) ([]byte, error) {
	return os.ReadFile("testdata/" + dir + "/" + file)
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
