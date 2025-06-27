package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"maps"
	"os"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
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

func buildUI() {
	a := app.New()
	w := a.NewWindow("Idle Planet Calculator")

	orderContainer := container.NewVBox()
	newOrderButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		item := NewOrder(items)
		orderContainer.Add(item)
		orders = append(orders, item)
		item.SetOnRemoved(func() {
			orderContainer.Remove(item)
			index := slices.Index(orders, item)
			orders = slices.Delete(orders, index, index+1)
		})
		item.SetOnItemChanged(func(name string) {
			gameItem := data[name]
			item.orderItem = gameItem
		})
	})

	getSeparator := func() *canvas.LinearGradient {
		return canvas.NewHorizontalGradient(color.RGBA{100, 100, 100, 255}, color.RGBA{22, 22, 22, 255})
	}

	resultTable := widget.NewTable(
		func() (rows int, cols int) {
			return len(results), 3
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("template")
			label.SizeName = theme.SizeNameCaptionText
			return label
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			var text string
			switch tci.Col {
			case 0:
				text = results[tci.Row].Item.Name
			case 1:
				text = humanize.Comma(int64(results[tci.Row].Amount))
			case 2:
				text = humanize.Comma(int64(results[tci.Row].Value))
			default:
				text = "Template"
			}
			co.(*widget.Label).SetText(text)
		},
	)
	resultTable.SetColumnWidth(0, 130)
	resultTable.SetColumnWidth(1, 100)
	resultTable.SetColumnWidth(2, 140)
	resultTable.ShowHeaderRow = true
	resultTable.CreateHeader = func() fyne.CanvasObject {
		label := widget.NewLabel("template")
		label.SizeName = theme.SizeNameCaptionText
		return label
	}
	calcHeader := func(id widget.TableCellID) string {
		switch id.Col {
		case 0:
			return "Item"
		case 1:
			return "Amount"
		case 2:
			return "Value"
		default:
			return "Template"
		}
	}
	resultTable.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		template.(*widget.Label).SetText(calcHeader(id))
	}

	displayResults := func(ingredients []Ingredient) {
		results = ingredients
		resultTable.Refresh()
	}

	sortResults := func(ingredients map[string]Ingredient) []Ingredient {
		result := slices.Collect(maps.Values(ingredients))
		slices.SortFunc(result, func(a, b Ingredient) int {
			if a.Item.Type != b.Item.Type {
				if a.Item.Type > b.Item.Type {
					return -1
				} else {
					return 1
				}
			}
			if a.Value != b.Value {
				if a.Value < b.Value {
					return 1
				} else {
					return -1
				}
			}
			return 0
		})
		return result
	}

	ordersLabel := widget.NewLabel("Orders")
	ordersLabel.SizeName = theme.SizeNameSubHeadingText
	newOrdersContainer := container.NewHBox(ordersLabel, layout.NewSpacer(), newOrderButton)
	newOrderButton.Resize(newOrderButton.MinSize())
	calculateButton := widget.NewButtonWithIcon("Calculate", theme.SettingsIcon(), func() {
		order := make([]Ingredient, 0)
		for _, o := range orders {
			order = append(order, Ingredient{
				Item:   o.orderItem,
				Amount: o.amount,
			})
		}
		result := calculateIngredients(order)
		displayResults(sortResults(result))
	})

	w.SetContent(
		container.NewBorder(
			container.NewVBox(
				newOrdersContainer,
				getSeparator(),
				orderContainer,
				calculateButton,
				getSeparator(),
			),
			nil,
			nil,
			nil,
			resultTable,
		))
	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}

func calculateIngredients(order []Ingredient) (bill map[string]Ingredient) {
	bill = make(map[string]Ingredient)
	for _, o := range order {
		ingredients := getIngredients(o.Item)
		for _, i := range ingredients {
			value := i.Item.Value * i.Amount * o.Amount
			amount := i.Amount * o.Amount
			if item, found := bill[i.Item.Name]; !found {
				bill[i.Item.Name] = Ingredient{i.Item, value, amount}
			} else {
				item.Value += value
				item.Amount += amount
			}
		}
	}
	return
}

func getIngredients(item GameItem) (ingredients map[string]Ingredient) {
	ingredients = make(map[string]Ingredient, 0)

	for _, i := range item.Ingredients {
		if newItem, found := ingredients[i.Item.Name]; !found {
			ingredients[i.Item.Name] = Ingredient{i.Item, i.Value, i.Amount}
		} else {
			newItem.Amount += i.Amount
		}

		if len(i.Item.Ingredients) > 0 {
			subIngredients := getIngredients(i.Item)
			for k, v := range subIngredients {
				if newItem, found := ingredients[k]; !found {
					ingredients[k] = Ingredient{v.Item, v.Value * i.Amount, v.Amount * i.Amount}
				} else {
					newItem.Amount += v.Amount * i.Amount
					ingredients[k] = newItem
				}
			}
		}
	}

	return
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
