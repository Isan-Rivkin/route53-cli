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

	for rowIdx := 0; rowIdx < 20; rowIdx++ {
		
		for colIdx, c := range cols {
			table.AddRow(rowIdx+1, colIdx, fmt.Sprintf("val=%s,%d,%d", c, rowIdx, colIdx))
		}
	}

	tableItem := table.Render(func(tsr *TableSelectionResult) {
		log.Info("Selected Item got, exiting!", *tsr)
		a.app.Stop()
	})

	// render grid
	grid := tview.NewGrid().
		SetRows(10, 0).
		//SetColumns(30, 0).
		SetBorders(true).
		AddItem(header.Render(), 0, 0, 1, 3, 0, 0, false)

	grid.AddItem(tableItem, 1, 0, 1, 3, 0, 0, true)

	return grid
}

func (a *DefaultApp) Render() error {

	appGrid := a.renderSimpleGrid()

	if err := a.app.SetRoot(appGrid, true).Run(); err != nil {
		log.WithError(err).Error("could not render table")
	}

	return nil
}
