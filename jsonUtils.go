package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func readJson(filePath string) []token {
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	tokens := []token{}
	json.Unmarshal([]byte(byteValue), &tokens)
	return tokens
}

func writeToJsonFile(data []token, fileName string) {
	file, _ := json.MarshalIndent(data, "", " ")
	_ = os.WriteFile(fileName, file, 0644)
}
