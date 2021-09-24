package cliui

import (
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"

	"github.com/rivo/tview"
)

type UIEventType string

const (
	R53TableSelection       UIEventType = "r53_table_selection"
	headerGridRowIdx        int         = 0
	r53GridTableRowIdx      int         = 1
	expansionViewGridRowIdx int         = 2
)

type AppEvent struct {
	Type         UIEventType
	EventPayload interface{}
}

type InteractiveApp interface {
	SetR53RecordsQueryResult(result *awsUtils.GetRecordAliasesResult)
	AddResourceExpansionTree(rootTxt string, childrenTxt []string)
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

func (app *R53App) AddResourceExpansionTree(rootTxt string, childrenTxt []string) {
	// rootGrid.RawGrid().AddItem(abstracts.NewInteractiveTree().RawTree, 2, 0, 1, 3, 0, 0, true)
	tree := abstracts.NewInteractiveTree().RawTree
	input := abstracts.NewDefaultGridItem(tree)
	input.Metadata.Row = expansionViewGridRowIdx
	input.Metadata.Focus = true
	app.App.AppendGridItem(input)
}

func (app *R53App) initDefaultRootGrid(header *tview.TextView, table *tview.Table) {
	rootGrid := abstracts.NewRootGrid()
	rootGrid.RawGrid().
		SetRows(10, 0).
		//SetColumns(30, 0).
		SetBorders(true)

	// add head
	headerInput := abstracts.NewDefaultGridItem(header)
	// header, headerGridRowIdx, 0, 1, 3, 0, 0, false
	headerInput.Metadata.Row = headerGridRowIdx
	rootGrid.RawGrid().AddItem(headerInput.ConvertToAddItemInput())

	// add table
	tableInput := abstracts.NewDefaultGridItem(table)
	tableInput.Metadata.Row = r53GridTableRowIdx
	tableInput.Metadata.Focus = true
	rootGrid.RawGrid().AddItem(tableInput.ConvertToAddItemInput())

	app.App.SetRootGrid(rootGrid.RawGrid())
}
