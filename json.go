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

func loadData() map[string]GameItem {
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

	gameItems := make(map[string]GameItem, 0)

	for _, ore := range data.Ores {
		gameItems[ore.Name] = GameItem{
			Name:  ore.Name,
			Value: ore.Value,
			Type:  Ore,
		}
	}

	for _, alloy := range data.Alloys {
		item := GameItem{
			Name:        alloy.Name,
			Type:        Alloy,
			Value:       alloy.Value,
			Ingredients: make([]Ingredient, 0),
		}

		for _, i := range alloy.Ingredients {
			ingredient, found := gameItems[i.Name]
			if !found {
				fmt.Printf("could not find: %s for %s\n", i.Name, alloy.Name)
				continue
			}
			item.Ingredients = append(item.Ingredients, Ingredient{
				Item:   ingredient,
				Amount: i.Amount,
			})
		}

		gameItems[alloy.Name] = item
	}

	for _, v := range data.Items {
		item := GameItem{
			Name:        v.Name,
			Type:        Item,
			Value:       v.Value,
			Ingredients: make([]Ingredient, 0),
		}

		for _, i := range v.Ingredients {
			ingredient, found := gameItems[i.Name]
			if !found {
				fmt.Printf("could not find: %s for %s", i.Name, v.Name)
				continue
			}
			item.Ingredients = append(item.Ingredients, Ingredient{
				Item:   ingredient,
				Amount: i.Amount,
			})
		}

		gameItems[v.Name] = item
	}

	return gameItems
}
