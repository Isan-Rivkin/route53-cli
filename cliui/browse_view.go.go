package cliui

import (
	"r53/cliui/abstracts"

	"github.com/rivo/tview"
)

func NewFocusTable(description *abstracts.TableDescription) *abstracts.TablePrompt {

	tablePrompt := abstracts.NewTable().PopulateTableDefault(description)
	return tablePrompt
}

type FocusView interface {
	Get() tview.Primitive
	ReRender(description *abstracts.TableDescription)
}

type FocusTableView struct {
	table *abstracts.TablePrompt
}

func (ft *FocusTableView) Get() tview.Primitive {
	return ft.table.Render()
}

func (ft *FocusTableView) ReRender(description *abstracts.TableDescription) {
	ft.table.CleanTable()
	ft.table.PopulateTableDefault(description)
}

func NewFocusView(description *abstracts.TableDescription) FocusView {
	return &FocusTableView{
		table: NewFocusTable(description),
	}
}
