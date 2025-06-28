package main

import (
	"strconv"

	"fyne.io/fyne/v2"
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
