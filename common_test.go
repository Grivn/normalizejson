package normalizejson

import (
	"encoding/json"
	"os"
)

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
