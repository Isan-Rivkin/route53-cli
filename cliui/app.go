package cliui

import (
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"

	"github.com/rivo/tview"
)

type UIEventType string

const (
	TableSelection UIEventType = "table_selection"
)

type AppEvent struct {
	Type         UIEventType
	EventPayload interface{}
}

type InteractiveApp interface {
	SetR53RecordsQueryResult(result *awsUtils.GetRecordAliasesResult)
	Run() error
}

type R53App struct {
	App              abstracts.App
	EventsController chan *AppEvent
}

func NewR53App(eventsController chan *AppEvent) InteractiveApp {
	return &R53App{
		EventsController: eventsController,
		App:              abstracts.NewApp(),
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
	// delete TEST
	here grid example also got channel sending update CHECKOUT the into sub command
	// rootGrid.RawGrid().AddItem(abstracts.NewInteractiveTree().RawTree, 2, 0, 1, 3, 0, 0, true)
	//
	app.App.SetRootGrid(rootGrid.RawGrid())
}
