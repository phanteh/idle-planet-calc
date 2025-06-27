package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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
	amount.OnChanged = func(s string) {
		val, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return
		}
		o.amount = int(val)
	}
	amount.Validator = validation.NewRegexp("^[0-9]+$", "Whole numbers only please")
	amount.SetText("1")

	itemSelector := widget.NewSelect(o.options, func(input string) {
		o.onItemChanged(input)
	})

	return &orderRenderer{
		order:        o,
		itemSelector: itemSelector,
		amount:       amount,
		remove: widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), func() {
			o.onRemoved()
		}),
	}
}

type orderRenderer struct {
	order        *Order
	itemSelector *widget.Select
	amount       *widget.Entry
	remove       *widget.Button
}

func (o *orderRenderer) Destroy() {
}

func (o *orderRenderer) Layout(size fyne.Size) {
	padding := float32(4)
	pos := fyne.NewPos(padding, padding)
	o.itemSelector.Move(pos)
	o.itemSelector.Resize(o.itemSelector.MinSize().AddWidthHeight(100, 0))
	pos.X += o.itemSelector.Size().Width + (padding * 2)
	o.amount.Move(pos)
	o.amount.Resize(o.amount.MinSize().AddWidthHeight(0, 0))
	pos.X = size.Width - (padding + o.remove.MinSize().Width)
	o.remove.Move(pos)
	o.remove.Resize(o.remove.MinSize())
}

func (o *orderRenderer) MinSize() fyne.Size {
	width := o.itemSelector.MinSize().Width +
		o.amount.Size().Width + 8
	height := o.itemSelector.MinSize().Height + 8
	return fyne.NewSize(width, height)
}

func (o *orderRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{o.itemSelector, o.amount, o.remove}
}

func (o *orderRenderer) Refresh() {
	o.itemSelector.Refresh()
	o.amount.Refresh()
}
