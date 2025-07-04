package main

import (
	"fmt"
	"image/color"
	"maps"
	"math"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
)

type App struct {
	app              fyne.App
	data             map[string]GameItem
	itemList         []string
	orders           []*Order
	results          []Ingredient
	craftMatBonus    float64
	smeltMatBonus    float64
	craftValBonus    float64
	smeltValBonus    float64
	underforgeBonus  float64
	dormsBonus       float64
	orderContainer   *fyne.Container
	bonusContainer   *fyne.Container
	resultSummary    *SummaryScreen
	summaryAccordion *widget.Accordion
	resultTable      *widget.Table
	craftMatEntry    *widget.Entry
	smeltMatEntry    *widget.Entry
	craftValEntry    *widget.Entry
	smeltValEntry    *widget.Entry
	underforgeEntry  *widget.Entry
	dormsEntry       *widget.Entry
}

func NewApp(data *jsonGameData) (app *App) {
	gameData := getGameData(data)
	app = &App{
		data:            gameData,
		itemList:        getItemList(gameData),
		orders:          make([]*Order, 0),
		results:         make([]Ingredient, 0),
		craftMatBonus:   1.0,
		smeltMatBonus:   1.0,
		craftValBonus:   1.0,
		smeltValBonus:   1.0,
		underforgeBonus: 1.0,
		dormsBonus:      1.0,
	}
	return
}

func (a *App) newOrderHandler() {
	item := NewOrder(a.itemList)
	a.orderContainer.Add(item)
	a.orders = append(a.orders, item)
	item.SetOnRemoved(func() {
		a.orderContainer.Remove(item)
		index := slices.Index(a.orders, item)
		a.orders = slices.Delete(a.orders, index, index+1)
	})
	item.SetOnItemChanged(func(name string) {
		gameItem := a.data[name]
		item.orderItem = gameItem
	})
}

func (a *App) getResultsTable() *widget.Table {
	resultTable := widget.NewTable(
		func() (rows int, cols int) {
			return len(a.results), 3
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
				text = a.results[tci.Row].Item.Name
			case 1:
				text = humanize.Comma(int64(a.results[tci.Row].Amount))
			case 2:
				text = fmt.Sprintf("$%s", humanize.Comma(int64(a.results[tci.Row].Value)))
			default:
				text = "Template"
			}
			co.(*widget.Label).SetText(text)
		},
	)
	resultTable.SetColumnWidth(0, 120)
	resultTable.SetColumnWidth(1, 80)
	resultTable.SetColumnWidth(2, 120)
	resultTable.ShowHeaderRow = true
	resultTable.CreateHeader = func() fyne.CanvasObject {
		label := widget.NewLabel("template")
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
	return resultTable
}

func (a *App) displayResults(ingredients []Ingredient) {
	a.results = ingredients
	a.resultTable.Refresh()
	a.resultSummary.Display(a.orders)
	a.resultSummary.Refresh()
	a.summaryAccordion.Refresh()
}

func (a *App) sortResults(ingredients map[string]Ingredient) []Ingredient {
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

func (a *App) calcResultsHandler() {
	order := make([]Ingredient, 0)
	for _, o := range a.orders {
		order = append(order, Ingredient{
			Item:   o.orderItem,
			Amount: o.amount,
		})
	}
	result := a.calculateIngredients(order)
	a.displayResults(a.sortResults(result))
}

func (a *App) getOrderAccordion(newOrderButton *widget.Button) *widget.Accordion {
	return widget.NewAccordion(
		widget.NewAccordionItem("Orders", container.NewVBox(
			a.orderContainer,
			container.NewHBox(newOrderButton),
			widget.NewAccordion(
				widget.NewAccordionItem("Bonuses", container.NewVBox(
					a.bonusContainer,
				)),
			)),
		),
	)
}

func (a *App) getBonuses() *fyne.Container {
	a.craftMatEntry = widget.NewEntry()
	a.smeltMatEntry = widget.NewEntry()
	a.craftValEntry = widget.NewEntry()
	a.smeltValEntry = widget.NewEntry()
	a.underforgeEntry = widget.NewEntry()
	a.dormsEntry = widget.NewEntry()

	getVal := func(input string) float64 {
		val, err := strconv.ParseFloat(input, 32)
		if err != nil {
			val = 1.0
		}
		return val
	}

	a.craftMatEntry.OnSubmitted = func(input string) {
		val := getVal(input)
		a.craftMatBonus = val
		a.craftMatEntry.SetText(fmt.Sprintf("%.2f", val))
	}
	a.smeltMatEntry.OnSubmitted = func(input string) {
		val := getVal(input)
		a.smeltMatBonus = val
		a.smeltMatEntry.SetText(fmt.Sprintf("%.2f", val))
	}
	a.craftValEntry.OnSubmitted = func(input string) {
		val := getVal(input)
		a.craftValBonus = val
		a.craftValEntry.SetText(fmt.Sprintf("%.2f", val))
	}
	a.smeltValEntry.OnSubmitted = func(input string) {
		val := getVal(input)
		a.smeltValBonus = val
		a.smeltValEntry.SetText(fmt.Sprintf("%.2f", val))
	}
	a.underforgeEntry.OnSubmitted = func(input string) {
		val := getVal(input)
		a.underforgeBonus = val
		a.underforgeEntry.SetText(fmt.Sprintf("%.2f", val))
	}
	a.dormsEntry.OnSubmitted = func(input string) {
		val := getVal(input)
		a.dormsBonus = val
		a.dormsEntry.SetText(fmt.Sprintf("%.2f", val))
	}

	a.craftMatEntry.SetText(fmt.Sprintf("%.2f", a.craftMatBonus))
	a.smeltMatEntry.SetText(fmt.Sprintf("%.2f", a.smeltMatBonus))
	a.craftValEntry.SetText(fmt.Sprintf("%.2f", a.craftValBonus))
	a.smeltValEntry.SetText(fmt.Sprintf("%.2f", a.smeltValBonus))
	a.underforgeEntry.SetText(fmt.Sprintf("%.2f", a.underforgeBonus))
	a.dormsEntry.SetText(fmt.Sprintf("%.2f", a.dormsBonus))

	materials := widget.NewForm(
		widget.NewFormItem("Smelt Material", a.smeltMatEntry),
		widget.NewFormItem("Craft Material", a.craftMatEntry),
		widget.NewFormItem("Underforge", a.underforgeEntry),
	)
	values := widget.NewForm(
		widget.NewFormItem("Smelt Value", a.smeltValEntry),
		widget.NewFormItem("Craft Value", a.craftValEntry),
		widget.NewFormItem("Dorms", a.dormsEntry),
	)
	return container.NewGridWithColumns(2, materials, values)
}

func (a *App) onStopped() {
	a.app.Preferences().SetFloat("craftMatBonus", a.craftMatBonus)
	a.app.Preferences().SetFloat("smeltMatBonus", a.smeltMatBonus)
	a.app.Preferences().SetFloat("dormsBonus", a.dormsBonus)
	a.app.Preferences().SetFloat("forgeBonus", a.underforgeBonus)
}

func (a *App) loadPreferences() {
	a.craftMatBonus = a.app.Preferences().FloatWithFallback("craftMatBonus", 1.0)
	a.smeltMatBonus = a.app.Preferences().FloatWithFallback("smeltMatBonus", 1.0)
	a.dormsBonus = a.app.Preferences().FloatWithFallback("dormsBonus", 1.0)
	a.underforgeBonus = a.app.Preferences().FloatWithFallback("forgeBonus", 1.0)
}

func (a *App) Run() {
	a.app = app.New()
	win := a.app.NewWindow("Idle Planet Calc")

	a.loadPreferences()
	a.orderContainer = container.NewVBox()
	a.bonusContainer = a.getBonuses()
	a.resultTable = a.getResultsTable()
	a.resultSummary = NewSummaryScreen()
	newOrderButton := widget.NewButtonWithIcon("Add ", theme.ContentAddIcon(), a.newOrderHandler)
	calculateButton := widget.NewButtonWithIcon("Calculate", theme.SettingsIcon(), a.calcResultsHandler)
	orderAccordion := a.getOrderAccordion(newOrderButton)
	a.summaryAccordion = widget.NewAccordion(widget.NewAccordionItem("Summary", a.resultSummary))
	a.app.Lifecycle().SetOnStopped(a.onStopped)

	win.SetContent(
		container.NewBorder(
			container.NewVBox(
				orderAccordion,
				calculateButton,
				getSeparator(),
				a.summaryAccordion,
			),
			nil,
			nil,
			nil,
			a.resultTable,
		))
	win.Resize(fyne.NewSize(400, 600))
	win.ShowAndRun()
}

func getSeparator() *canvas.LinearGradient {
	return canvas.NewHorizontalGradient(
		color.RGBA{100, 100, 100, 255},
		color.RGBA{22, 22, 22, 255})
}

func getItemList(data map[string]GameItem) (items []string) {
	items = make([]string, 0)

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

	for _, item := range sortedItem {
		if item.Type == Ore {
			continue
		}
		items = append(items, item.Name)
	}
	return
}

func getGameData(data *jsonGameData) map[string]GameItem {
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

func (a *App) calculateIngredients(order []Ingredient) (bill map[string]Ingredient) {
	bill = make(map[string]Ingredient)
	for _, o := range order {
		ingredients := a.getIngredients(o.Item)
		for _, i := range ingredients {
			value := i.Value * i.Amount * o.Amount
			amount := i.Amount * o.Amount
			if item, found := bill[i.Item.Name]; !found {
				bill[i.Item.Name] = Ingredient{i.Item, value, amount}
			} else {
				item.Value += value
				item.Amount += amount
				bill[i.Item.Name] = item
			}
		}
	}
	return
}

func (a *App) getBonusedMaterialAmount(itemType ItemType, value int) int {
	var roomBonus, projectBonus float64
	amount := float64(value)

	if itemType == Item {
		roomBonus = a.dormsBonus
		projectBonus = a.craftMatBonus
	} else {
		roomBonus = a.underforgeBonus
		projectBonus = a.smeltMatBonus
	}

	basePrice := amount - (amount * (roomBonus - 1))
	smeltBonus := basePrice * (projectBonus - 1)
	if smeltBonus < 1 {
		smeltBonus = math.Round(smeltBonus)
	}
	return int(math.Round(basePrice - smeltBonus))
}

func (a *App) getBonusedValue(itemType ItemType, value int) int {
	var projectBonus float64
	amount := float64(value)
	if itemType == Item {
		projectBonus = a.craftValBonus
	} else {
		projectBonus = a.smeltValBonus
	}
	return int(math.Round(amount * projectBonus))
}

func (a *App) getIngredients(item GameItem) (ingredients map[string]Ingredient) {
	ingredients = make(map[string]Ingredient, 0)
	for _, i := range item.Ingredients {
		amount := a.getBonusedMaterialAmount(item.Type, i.Amount)
		value := a.getBonusedValue(i.Item.Type, i.Item.Value)

		if newItem, found := ingredients[i.Item.Name]; !found {
			ingredients[i.Item.Name] = Ingredient{i.Item, value, amount}
		} else {
			newItem.Amount += amount
		}

		if len(i.Item.Ingredients) > 0 {
			subIngredients := a.getIngredients(i.Item)
			for k, v := range subIngredients {
				if newItem, found := ingredients[k]; !found {
					ingredients[k] = Ingredient{v.Item, v.Value, v.Amount * amount}
				} else {
					newItem.Amount += v.Amount * amount
					ingredients[k] = newItem
				}
			}
		}
	}
	return
}
