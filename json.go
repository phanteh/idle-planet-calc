package main

import (
	_ "embed"
	"encoding/json"
)

//go:embed inventory.json
var inventoryBytes []byte

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
	var data jsonGameData
	json.Unmarshal(inventoryBytes, &data)
	return &data
}
