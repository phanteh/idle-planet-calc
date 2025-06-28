package main

import (
	"image/color"
	"maps"
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

	accordion := widget.NewAccordion(
		widget.NewAccordionItem(
			"Orders",
			container.NewVBox(
				orderContainer,
				container.NewHBox(newOrderButton),
			)))

	w.SetContent(
		container.NewBorder(
			container.NewVBox(
				accordion,
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

type Order struct {
	widget.BaseWidget
	options       []string
	orderItem     GameItem
	amount        int
	onRemoved     func()
	onItemChanged func(string)
}

func NewOrder(options []string) *Order {
	item := &Order{
		options: options,
	}
	item.ExtendBaseWidget(item)
	return item
}

func (o *Order) SetOnItemChanged(onItemChanged func(string)) {
	o.onItemChanged = onItemChanged
}

func (o *Order) SetOnRemoved(onRemoved func()) {
	o.onRemoved = onRemoved
}

func (o *Order) CreateRenderer() fyne.WidgetRenderer {
	amount := widget.NewEntry()
	amount.SetText("1")
	o.amount = 1
	amount.OnChanged = func(s string) {
		val, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			o.amount = 1
			amount.SetText("1")
		}
		o.amount = int(val)
	}
	amount.Resize(amount.MinSize())

	decrement := widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), func() {
		o.amount -= 1
		amount.SetText(strconv.Itoa(o.amount))
	})
	increment := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), func() {
		o.amount += 1
		amount.SetText(strconv.Itoa(o.amount))
	})
	remove := widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
		o.onRemoved()
	})

	decrement.Resize(decrement.MinSize().AddWidthHeight(-10, -10))
	increment.Resize(increment.MinSize().AddWidthHeight(-10, -10))
	remove.Resize(remove.MinSize().AddWidthHeight(-10, -10))

	itemSelector := widget.NewSelect(o.options, func(input string) {
		o.onItemChanged(input)
	})
	itemSelector.Resize(itemSelector.MinSize().AddWidthHeight(50, 0))

	return &orderRenderer{
		order:        o,
		itemSelector: itemSelector,
		amount:       amount,
		increment:    increment,
		decrement:    decrement,
		remove:       remove,
	}
}

type orderRenderer struct {
	order                        *Order
	itemSelector                 *widget.Select
	amount                       *widget.Entry
	increment, decrement, remove *widget.Button
}

func (o *orderRenderer) Destroy() {
}

func (o *orderRenderer) Layout(size fyne.Size) {
	padding := float32(4)
	pos := fyne.NewPos(padding, padding)
	buttonHeight := (size.Height / 2) - (o.increment.Size().Height / 2)
	o.itemSelector.Move(pos)
	pos.X += o.itemSelector.Size().Width + (padding * 2)
	pos.Y = buttonHeight
	o.decrement.Move(pos)
	pos.X += o.decrement.Size().Width + (padding * 2)
	pos.Y = padding
	o.amount.Move(pos)
	pos.X += o.amount.Size().Width + (padding * 2)
	pos.Y = buttonHeight
	o.increment.Move(pos)
	pos.X += o.increment.Size().Width + (padding * 2)
	pos.X = size.Width - (padding + o.remove.MinSize().Width)
	o.remove.Move(pos)
}

func (o *orderRenderer) MinSize() fyne.Size {
	width := o.itemSelector.MinSize().Width +
		o.amount.Size().Width + 8
	height := o.itemSelector.MinSize().Height + 8
	return fyne.NewSize(width, height)
}

func (o *orderRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{o.itemSelector, o.amount, o.increment, o.decrement, o.remove}
}

func (o *orderRenderer) Refresh() {
	o.itemSelector.Refresh()
	o.amount.Refresh()
}
