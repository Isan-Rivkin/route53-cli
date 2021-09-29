package v1

import (
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"

	"github.com/gdamore/tcell/v2"
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
	UpdateFocusTable(tableType awsUtils.ResourceType, resources []awsUtils.Resource) error
	UpdateFocusTxtView(rType awsUtils.ResourceType, resources []awsUtils.Resource) error
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

func (app *R53App) UpdateFocusTxtView(rType awsUtils.ResourceType, resources []awsUtils.Resource) error {
	txtView, err := app.RenderFocusTxtViewItems(rType, resources)

	if err != nil {
		return err
	}
	input := abstracts.NewDefaultGridItem(txtView.Render())
	input.Metadata.Row = expansionViewGridRowIdx
	input.Metadata.Column = 1
	input.Metadata.Focus = true

	app.App.AppendGridItem(input)

	return nil
}

func (app *R53App) UpdateFocusTable(tableType awsUtils.ResourceType, resources []awsUtils.Resource) error {
	table, err := app.RenderFocusTable(tableType, resources)

	if err != nil {
		return err
	}

	input := abstracts.NewDefaultGridItem(table)
	input.Metadata.Row = expansionViewGridRowIdx
	input.Metadata.Column = 1
	input.Metadata.Focus = true

	app.App.AppendGridItem(input)

	return nil
}

func (app *R53App) AddResourceExpansionTree(rootTxt string, childrenTxt []string) {
	// rootGrid.RawGrid().AddItem(abstracts.NewInteractiveTree().RawTree, 2, 0, 1, 3, 0, 0, true)
	var children []*NodeInput
	rootNode := &NodeInput{Txt: rootTxt, Color: tcell.ColorGreen, Selectable: true, Expanded: true}
	for _, c := range childrenTxt {
		n := &NodeInput{Txt: c, Selectable: true, Expanded: false, Color: tcell.ColorYellow}
		children = append(children, n)
	}
	root := NewRootNodeWithChildren(rootNode, children)
	tree := abstracts.NewInteractiveTree(root, func(node *tview.TreeNode) {
		node.SetColor(tcell.ColorBlue)
	}).RawTree
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
