package abstracts

import (
	"fmt"

	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

type AppOptions struct {
	EnableMouse bool
}

func DefaultOpts() *AppOptions {
	return &AppOptions{
		EnableMouse: true,
	}
}

type App interface {
	SetRootGrid(rootGrid *tview.Grid)
	Render() error
}

type DefaultApp struct {
	app *tview.Application
}

func NewApp() App {
	opts := DefaultOpts()
	return &DefaultApp{
		app: tview.NewApplication().EnableMouse(opts.EnableMouse),
	}
}

func getDummyTable() *tview.Table {
	// parse table
	table := NewTable()
	cols := []string{"AAA", "BNBB", "CCC", "DDDDD"}
	table.AddHeaders(cols)

	for rowIdx := 0; rowIdx < 20; rowIdx++ {

		for colIdx, c := range cols {
			table.AddRow(rowIdx+1, colIdx, fmt.Sprintf("dummy=%s,%d,%d", c, rowIdx, colIdx))
		}
	}

	table.AddSelectionCallBack(func(tsr *TableSelectionResult) {
		log.Info("Selected Item got, exiting!", *tsr)
	})

	tableItem := table.Render()
	return tableItem
}
func (a *DefaultApp) renderSimpleGrid() *tview.Grid {
	//appGrid := tview.NewGrid().SetBorders(true) //.AddItem(tp.table, 0, 0, 1, 3, 0, 0, true)
	header := NewLabeledText(&LabeledTextViewInput{
		Rows: []*LabeledTextRowInput{
			{
				LabelTxt:   "Hosted Zone:",
				LabelColor: OrangeTxtColor,
				ValColor:   WhiteTxtColor,
				ValTxt:     "api.company.com",
			},
			{
				LabelTxt:   "ID:",
				LabelColor: OrangeTxtColor,
				ValTxt:     "/hostedzone/1237LZW123Z2BH",
				ValColor:   WhiteTxtColor,
			},
			{
				LabelTxt:   "Total Records:",
				LabelColor: OrangeTxtColor,
				ValTxt:     "534",
				ValColor:   WhiteTxtColor,
			},
			{
				LabelTxt:   "Private:",
				LabelColor: OrangeTxtColor,
				ValTxt:     "true",
				ValColor:   WhiteTxtColor,
			},
			{
				LabelTxt:   "NS Verified: ",
				LabelColor: OrangeTxtColor,
				ValColor:   GreenTxtColor,
				ValTxt:     "true",
			},
		},
	})
	// parse table
	table := NewTable()
	cols := []string{"Col1", "Col2", "Col3"}
	table.AddHeaders(cols)

	for rowIdx := 0; rowIdx < 250; rowIdx++ {

		for colIdx, c := range cols {
			table.AddRow(rowIdx+1, colIdx, fmt.Sprintf("val=%s,%d,%d", c, rowIdx, colIdx))
		}
	}
	table.AddSelectionCallBack(func(tsr *TableSelectionResult) {
		log.Info("Selected Item got, exiting!", *tsr)
		a.app.Stop()
	})
	tableItem := table.Render()

	// render grid
	grid := tview.NewGrid().
		SetRows(10, 0).
		//SetColumns(30, 0).
		SetBorders(true).
		AddItem(header.Render(), 0, 0, 1, 3, 0, 0, false)

	grid.AddItem(tableItem, 1, 0, 1, 3, 0, 0, true)

	return grid
}

func (a *DefaultApp) SetRootGrid(rootGrid *tview.Grid) {
	a.app.SetRoot(rootGrid, true)
}

func (a *DefaultApp) Render() error {

	//appGrid := a.renderSimpleGrid()

	// example: update dynamically the ui by expanding the grid and adding another table
	// go func() {
	// 	time.Sleep(time.Second * 3)
	// 	a.app.QueueUpdateDraw(func() {
	// 		appGrid.SetRows(10, 0, 0)
	// 		t := getDummyTable()
	// 		appGrid.AddItem(t, 2, 0, 1, 3, 0, 0, true)
	// 	})
	// }()

	if err := a.app.Run(); err != nil {
		log.WithError(err).Error("could not render table")
	}

	return nil
}
