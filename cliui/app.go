package cliui

import (
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"

	"github.com/rivo/tview"
)

type InteractiveApp interface {
	SetR53RecordsQueryResult(result *awsUtils.GetRecordAliasesResult)
	Run() error
}

type R53App struct {
	App abstracts.App
}

func NewR53App() InteractiveApp {
	return &R53App{
		App: abstracts.NewApp(),
	}
}

func (app *R53App) Run() error {
	return app.App.Render()
}

func (app *R53App) initDefaultRootGrid(header *tview.TextView, table *tview.Table) {
	rootGrid := abstracts.NewRootGrid()
	rootGrid.RawGrid().
		SetRows(10, 0).
		//SetColumns(30, 0).
		SetBorders(true).
		AddItem(header, 0, 0, 1, 3, 0, 0, false)

	rootGrid.RawGrid().AddItem(table, 1, 0, 1, 3, 0, 0, true)
	app.App.SetRootGrid(rootGrid.RawGrid())
}
