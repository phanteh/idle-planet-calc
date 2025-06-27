package main

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
