package abstracts

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TableSelectionResult struct {
	IsSelected   bool
	SelectedText string
	RowSelected  int
	ColSelected  int
}

type OnTableSelectionFunc func(*TableSelectionResult)

type TablePrompt struct {
	table *tview.Table
}

// NewTable create new table view
func NewTable() *TablePrompt {

	return &TablePrompt{
		table: tview.NewTable().SetBorders(false),
	}
}

// AddHeaders add headers to the grid
func (tp *TablePrompt) AddHeaders(headers []string) {

	for i, headerText := range headers {
		tp.table.SetCell(0, i, tview.NewTableCell(strings.ToUpper(headerText)).
			SetSelectable(false).
			SetTextColor(tcell.ColorWhite).
			SetTransparency(true).
			SetAttributes(tcell.AttrBold).SetAlign(tview.AlignLeft))
	}

}

// AddRow adding new row to the grid
func (tp *TablePrompt) AddRow(row, column int, text string) {
	tp.table.SetCell(row, column, tview.NewTableCell(fmt.Sprintf("%s    ", text)).SetTextColor(tcell.ColorGreen).SetAlign(tview.AlignLeft))
}

// Render returns the table that needs to be show on the screen
func (tp *TablePrompt) Render(onTableSelected OnTableSelectionFunc) *tview.Table {

	tp.table.SetSelectable(true, false)
	tp.table.SetSelectedFunc(func(row int, column int) {
		tp.table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		tp.table.SetSelectable(true, false)
		//tp.app.Stop()
		cellVall := tp.table.GetCell(row, column).Text
		onTableSelected(&TableSelectionResult{
			IsSelected:   true,
			SelectedText: cellVall,
			RowSelected:  row,
			ColSelected:  column,
		})
	})

	return tp.table
}
