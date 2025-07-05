package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
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

type SummaryScreen struct {
	widget.BaseWidget
	ingredients []*ResultSummary
	renderer    *summaryScreenRenderer
}

func NewSummaryScreen() *SummaryScreen {
	item := &SummaryScreen{
		ingredients: make([]*ResultSummary, 0),
	}
	item.ExtendBaseWidget(item)
	return item
}

func (s *SummaryScreen) Display(ingredients []ResultItem) {
	s.ingredients = make([]*ResultSummary, 0)
	for _, i := range ingredients {
		s.ingredients = append(s.ingredients, NewResultSummary(i))
	}
	s.renderer.container = container.NewVBox()
}

func (s *SummaryScreen) CreateRenderer() fyne.WidgetRenderer {
	item := &summaryScreenRenderer{
		summaryScreen: s,
		container:     container.NewVBox(),
	}
	s.renderer = item
	return item
}

type summaryScreenRenderer struct {
	summaryScreen *SummaryScreen
	container     *fyne.Container
}

func (s *summaryScreenRenderer) Destroy() {
}

func (s *summaryScreenRenderer) Layout(size fyne.Size) {
	s.container.Resize(size)
	s.CheckChildren()
	s.CheckLabels()
}

func (s *summaryScreenRenderer) CheckLabels() {
	minWidth := float32(0)
	for _, o := range s.container.Objects {
		if item, ok := o.(*ResultSummary); ok {
			width := item.renderer.itemLabel.MinSize().Width
			if item.renderer.valueLabel.MinSize().Width > width {
				width = item.renderer.valueLabel.MinSize().Width
			}
			if minWidth < width {
				minWidth = width
			}
		}
	}

	for _, o := range s.container.Objects {
		if item, ok := o.(*ResultSummary); ok {
			item.setLabelWidth(minWidth)
		}
	}
}

func (s *summaryScreenRenderer) MinSize() fyne.Size {
	return s.container.MinSize()
}

func (s *summaryScreenRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{s.container}
}

func (s *summaryScreenRenderer) CheckChildren() {
	if len(s.container.Objects) == 0 &&
		len(s.summaryScreen.ingredients) > 0 {
		total := 0
		for index, ingredient := range s.summaryScreen.ingredients {
			total += ingredient.result.Value
			if index > 0 {
				s.container.Add(getSeparator())
			}
			s.container.Add(ingredient)
		}
		s.container.Add(getSeparator())
		s.container.Add(widget.NewLabel(
			fmt.Sprintf("Total: $%s", humanize.Comma(int64(total)))))
	}
}

func (s *summaryScreenRenderer) Refresh() {
	s.CheckChildren()
	s.CheckLabels()
	s.container.Refresh()
}

type ResultSummary struct {
	widget.BaseWidget
	labelWidth float32
	result     ResultItem
	renderer   *resultSummaryRenderer
}

func NewResultSummary(result ResultItem) *ResultSummary {
	item := &ResultSummary{
		result: result,
	}
	item.ExtendBaseWidget(item)
	return item
}

func (r *ResultSummary) setLabelWidth(value float32) {
	r.labelWidth = value
}

func (r *ResultSummary) CreateRenderer() fyne.WidgetRenderer {
	nameLabel := widget.NewLabel(fmt.Sprintf("%d x %s",
		r.result.Amount, r.result.Name))
	valueLabel := widget.NewLabel(fmt.Sprintf("$%s",
		humanize.Comma(int64(r.result.Value))))
	valueLabel.SizeName = theme.SizeNameCaptionText

	subContainer := container.NewVBox()

	renderer := &resultSummaryRenderer{
		resultSummary: r,
		itemLabel:     nameLabel,
		valueLabel:    valueLabel,
		subContainer:  subContainer,
	}

	r.renderer = renderer
	return renderer
}

type resultSummaryRenderer struct {
	resultSummary *ResultSummary
	itemLabel     *widget.Label
	valueLabel    *widget.Label
	subContainer  *fyne.Container
}

func (r *resultSummaryRenderer) Destroy() {
}

func (r *resultSummaryRenderer) Layout(fyne.Size) {
	r.CheckChildren()
	padding := float32(2)
	pos := fyne.NewPos(padding, padding)

	r.itemLabel.Move(pos)
	r.itemLabel.Resize(fyne.NewSize(r.resultSummary.labelWidth, r.itemLabel.MinSize().Height))

	widest := r.itemLabel.Size().Width
	if r.valueLabel.Size().Width > widest {
		widest = r.valueLabel.Size().Width
	}

	pos.X = widest + (padding * 4)
	r.subContainer.Resize(r.subContainer.MinSize())
	r.subContainer.Move(pos)

	pos.X = padding
	pos.Y = r.itemLabel.Size().Height + padding
	r.valueLabel.Resize(r.valueLabel.MinSize())
	r.valueLabel.Move(pos)
}

func (r *resultSummaryRenderer) MinSize() fyne.Size {
	size := r.itemLabel.MinSize()

	if r.valueLabel.MinSize().Width > size.Width {
		size.Width = r.valueLabel.MinSize().Width
	}

	if r.resultSummary.labelWidth > size.Width {
		size.Width = float32(r.resultSummary.labelWidth)
	}

	size.Width += r.subContainer.MinSize().Width
	size.Height += r.valueLabel.MinSize().Height

	if r.subContainer.MinSize().Height > size.Height {
		size.Height = r.subContainer.MinSize().Height
	}

	return size
}

func (r *resultSummaryRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.itemLabel, r.valueLabel, r.subContainer}
}

func (r *resultSummaryRenderer) CheckChildren() {
	if len(r.subContainer.Objects) == 0 &&
		len(r.resultSummary.result.Ingredients) > 0 {
		for index, ingredient := range r.resultSummary.result.Ingredients {
			if index > 0 {
				r.subContainer.Add(getSeparator())
			}
			label := widget.NewLabel(fmt.Sprintf("%d x %s", ingredient.Amount, ingredient.Name))
			label.SizeName = theme.SizeNameCaptionText
			r.subContainer.Add(label)
		}
	}
}

func (r *resultSummaryRenderer) Refresh() {
	r.CheckChildren()
	r.subContainer.Refresh()
}
