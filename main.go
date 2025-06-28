package main

import (
	"maps"
	"slices"
)

var (
	data    map[string]GameItem
	items   []string
	orders  []*Order
	results []Ingredient
)

func main() {
	data = loadData()
	sortedItem := slices.Collect(maps.Values(data))
	slices.SortFunc(sortedItem, func(a, b GameItem) int {
		if a.Type > b.Type {
			return -1
		}
		if a.Type < b.Type {
			return 1
		}
		if a.Value == b.Value {
			return 0
		}
		if a.Value > b.Value {
			return -1
		}
		return 1
	})
	items = make([]string, 0)
	for _, item := range sortedItem {
		if item.Type == Ore {
			continue
		}
		items = append(items, item.Name)
	}
	buildUI()
}
