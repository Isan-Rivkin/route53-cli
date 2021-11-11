package cliui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ExpansionIndicatorManager struct {
	path         []string
	onPathChange func(path []string, back bool)
}

func NewExpansionPathIndicator(initial string, onPathChange func(path []string, back bool)) *ExpansionIndicatorManager {
	path := []string{}

	if initial != "" {
		path = append(path, initial)
	}

	return &ExpansionIndicatorManager{
		onPathChange: onPathChange,
		path:         path,
	}

}

func (e *ExpansionIndicatorManager) UpdatePath(newPath string, back bool) {
	if back && len(e.path) > 0 {
		e.path = e.path[:len(e.path)-1]
	}
	e.path = append(e.path, newPath)

	e.onPathChange(e.path, back)
}

type IndicatorItem interface {
	Get() tview.Primitive
	MakeCurrentIndicator()
	UnmakeCurrentIndicator()
}

type IndicatorText struct {
	item *tview.TextView
}

func NewIndicatorItem(name string, current bool) IndicatorItem {

	t := tview.
		NewTextView().
		SetText(fmt.Sprintf("< %s >", name))
	t.SetTextAlign(tview.AlignCenter)

	indicator := &IndicatorText{item: t}

	if current {
		indicator.MakeCurrentIndicator()
	} else {
		indicator.UnmakeCurrentIndicator()
	}

	return indicator
}

func (ii *IndicatorText) Get() tview.Primitive {
	return ii.item
}

func (ii *IndicatorText) MakeCurrentIndicator() {
	ii.item.SetBackgroundColor(tcell.ColorDarkOrange)
	ii.item.SetTextColor(tcell.ColorWhite)
}
func (ii *IndicatorText) UnmakeCurrentIndicator() {
	ii.item.SetBackgroundColor(tcell.ColorDarkBlue)
	ii.item.SetTextColor(tcell.ColorWhite)
}

type IndicatorWidget struct {
	path []IndicatorItem
	flex *tview.Flex
}

func NewIndicatorWidget(initial string) *IndicatorWidget {
	path := []IndicatorItem{}
	flex := tview.NewFlex().SetDirection(tview.FlexColumn)
	flex.SetBorderPadding(1, 1, 0, 0)
	if initial != "" {
		path = append(path, NewIndicatorItem(initial, true))
	}

	return &IndicatorWidget{
		path: path,
		flex: flex,
	}
}

func (iw *IndicatorWidget) OnPathChange(newPath string, back bool) {

	for _, w := range iw.path {
		w.UnmakeCurrentIndicator()
	}

	if back {
		if len(iw.path) > 0 {
			lastIdx := len(iw.path) - 1
			lastItem := iw.path[lastIdx]
			iw.path = iw.path[:lastIdx]
			iw.flex.RemoveItem(lastItem.Get())
		}
	} else {
		item := NewIndicatorItem(newPath, true)
		item.MakeCurrentIndicator()
		iw.path = append(iw.path, item)
		//iw.flex.AddItem(item.Get(), 0, 1, false)
		iw.addListItem(item)
	}

}

func (iw *IndicatorWidget) addListItem(item IndicatorItem) {
	iw.flex.AddItem(item.Get(), 20, 1, false)
}
func (iw *IndicatorWidget) GetRender() tview.Primitive {

	for _, p := range iw.path {
		//iw.flex.AddItem(p.Get(), 0, 1, false)
		iw.addListItem(p)
	}

	return iw.flex
}

func NewIndicatorView(initial string) *IndicatorWidget {
	return NewIndicatorWidget(initial)
}
