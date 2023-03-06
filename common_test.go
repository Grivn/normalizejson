package gojson

import "os"

func readTestData(dir, file string) ([]byte, error) {
	return os.ReadFile("testdata/" + dir + "/" + file)
}
