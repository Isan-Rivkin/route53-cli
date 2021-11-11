package abstracts

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TableDescription struct {
	Headers                  []string
	Rows                     []map[string]string
	RowReference             []interface{}
	OnTableSelectionCallback func(*TableSelectionResult)
}

type TableSelectionResult struct {
	IsSelected        bool
	SelectedText      string
	SelectedReference interface{}
	RowSelected       int
	ColSelected       int
	RowCells          []*tview.TableCell
}

type OnTableSelectionFunc func(*TableSelectionResult)

type TablePrompt struct {
	table           *tview.Table
	onTableSelected OnTableSelectionFunc
}

// NewTable create new table view
func NewTable() *TablePrompt {

	return &TablePrompt{
		table: tview.NewTable().SetBorders(false),
	}
}

// Add call back functions
func (tp *TablePrompt) AddSelectionCallBack(onTableSelected OnTableSelectionFunc) {
	tp.onTableSelected = onTableSelected
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

func (tp *TablePrompt) defaultCell(text string, color tcell.Color) *tview.TableCell {
	cell := tview.NewTableCell(fmt.Sprintf("%s    ", text)).SetTextColor(color).SetAlign(tview.AlignLeft)
	return cell
}

// AddRow adding new row to the grid
// reference is some identifier on top of the actual cell value such as some host id etc
func (tp *TablePrompt) AddRow(row, column int, text string, reference interface{}) {
	cell := tp.defaultCell(text, tcell.ColorWhite)
	cell.SetReference(reference)
	tp.table.SetCell(row, column, cell)
}

func (tp *TablePrompt) RemoveRow(row int) {
	tp.table.RemoveRow(row)
}

func (tp *TablePrompt) CleanTable() {
	tp.table.Clear()
}

func (tp *TablePrompt) PopulateTableDefault(description *TableDescription) *TablePrompt {

	tp.AddSelectionCallBack(func(tsr *TableSelectionResult) {
		if description.OnTableSelectionCallback != nil {
			description.OnTableSelectionCallback(tsr)
		}
	})

	header := description.Headers

	tp.AddHeaders(append([]string{"#"}, header...))

	currentRow := 0

	var rowRefValue interface{}

	// populate columns
	for idx, row := range description.Rows {

		if len(description.RowReference) > idx {
			rowRefValue = description.RowReference[idx]
		}

		// add the # col for row number
		// populate cell
		tp.AddRow(currentRow+1, 0, fmt.Sprintf("%d", currentRow+1), rowRefValue)

		// add all the cols in the result
		for headerNum, headerName := range header {
			// populate cell
			tp.AddRow(currentRow+1, headerNum+1, row[headerName], rowRefValue)
		}
		currentRow++
	}

	return tp
}

func (tp *TablePrompt) AddInfoRow(row, column int, text string) {
	cell := tp.defaultCell(text, tcell.ColorGreen).SetExpansion(100)
	tp.table.SetCell(row, column, cell)
}

func (tp *TablePrompt) onSelected(row int, column int) {
	tp.table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	tp.table.SetSelectable(true, false)
	selectedCell := tp.table.GetCell(row, column)
	cellVall := selectedCell.Text
	if tp.onTableSelected != nil {
		// get all cells
		rowCells := []*tview.TableCell{}
		colsNum := tp.table.GetColumnCount()
		for i := 0; i < colsNum; i++ {
			val := tp.table.GetCell(row, i)
			rowCells = append(rowCells, val)
		}
		tp.onTableSelected(&TableSelectionResult{
			IsSelected:        true,
			SelectedText:      cellVall,
			SelectedReference: selectedCell.Reference,
			RowSelected:       row,
			ColSelected:       column,
			RowCells:          rowCells,
		})
	}
}

// Render returns the table that needs to be show on the screen
func (tp *TablePrompt) Render() *tview.Table {

	tp.table.SetSelectable(true, false)
	tp.table.SetSelectedFunc(tp.onSelected)

	return tp.table
}
