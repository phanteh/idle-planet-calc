package main

type ItemType int

const (
	Ore ItemType = iota
	Alloy
	Item
)

var itemTypeName = map[ItemType]string{
	Ore:   "Ore",
	Alloy: "Alloy",
	Item:  "Item",
}

func (it ItemType) String() string {
	return itemTypeName[it]
}

type GameItem struct {
	Name        string
	Type        ItemType
	Value       int
	Ingredients []Ingredient
}

type Ingredient struct {
	Item   GameItem
	Value  int
	Amount int
}
