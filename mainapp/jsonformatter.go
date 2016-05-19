package main

import (
	"encoding/json"
	"io/ioutil"
)

func generateJson(events []Event) ([]byte, error) {
	b, err := json.MarshalIndent(events, "", "    ")
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}

func saveJsonToFile(filename string, jsonStr []byte) error {
	err := ioutil.WriteFile(filename, jsonStr, 0644)
	return err
}
