package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type jsonGameData struct {
	Ores   []jsonGameItem `json:"ores"`
	Alloys []jsonGameItem `json:"alloys"`
	Items  []jsonGameItem `json:"items"`
}

type jsonGameItem struct {
	Name        string           `json:"name"`
	Value       int              `json:"value"`
	Ingredients []jsonIngredient `json:"ingredients"`
}

type jsonIngredient struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

func loadData() *jsonGameData {
	jsonFile, err := os.Open("inventory.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var data jsonGameData

	json.Unmarshal(bytes, &data)

	return &data
}
