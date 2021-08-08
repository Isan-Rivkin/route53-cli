package abstracts

import (
	"github.com/rivo/tview"
)

func newGridCell(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
}

type MultiLineTable interface {
	RawGrid() *tview.Grid
}

type TableGrid struct {
	rows    int
	cols    int
	matrix  [][]tview.Primitive
	headers []string
	Grid    *tview.Grid
}

func NewMultiLineTable(totalRows, totalCols int) MultiLineTable {

	m := make([][]tview.Primitive, totalRows)

	for i := 0; i < totalRows; i++ {
		m[i] = make([]tview.Primitive, totalCols)
	}

	return &TableGrid{
		Grid:   tview.NewGrid(),
		matrix: m,
	}
}

func (g *TableGrid) RawGrid() *tview.Grid {
	return g.Grid
}

func (g *TableGrid) Render() *tview.Grid {
	return g.Grid
}

func (g *TableGrid) AddCell(row, column int, text string) {
	g.matrix[row][column] = newGridCell(text)
}

func (tp *TableGrid) AddSelectionCallBack(onTableSelected OnTableSelectionFunc) {

}

func (tp *TableGrid) AddHeaders(headers []string) {
	tp.headers = headers
}
