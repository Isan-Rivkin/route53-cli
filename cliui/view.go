package cliui

import (
	"fmt"
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"

	"github.com/rivo/tview"
)

const (
	headerGridRowIdx  int = 0
	nextGridRowIdx    int = 1
	focusViewGridIdx  int = 2
	footerViewGridIdx int = 3
)

type View interface {
	Run()
}

type ViewController struct {
	app        *tview.Application
	startFocus tview.Primitive
	pages      *tview.Pages
	// indicates the resource expansion path of a user r53 -> alb -> ec2
	expandIndicator *ExpansionIndicatorManager
	// manages the table of resources
	focusView FocusView
}

func NewViewController() View {

	c := &ViewController{}
	c.init()
	return c
}

func getMockData() *abstracts.TableDescription {
	headers := []string{"Name", "Type", "Hosted", "DNS", "CPU"}
	rows := []map[string]string{}
	for i := 0; i < 16; i++ {
		r := map[string]string{}
		for j, h := range headers {
			r[h] = fmt.Sprintf("%d-%s-%d", i, h, j)
		}
		rows = append(rows, r)
	}

	desc := &abstracts.TableDescription{
		Headers: headers,
		Rows:    rows,
	}
	return desc
}

func (v *ViewController) init() {

	// init app
	v.app = tview.NewApplication()
	// init pages
	v.pages = tview.NewPages()

	// add browse screen
	browseScreen := v.createBrowseScreen(nil)
	v.pages.AddPage("browse", browseScreen, true, false)

	// add start view
	startScreen := v.createStartScreen(func(out *StartPageOutput) {
		// get resources from controller
		desc := getMockData()
		// update table
		v.focusView.ReRender(desc)

		go func() {
			v.app.QueueUpdateDraw(func() {
				// update indicator initial resource we started from
				v.expandIndicator.UpdatePath(out.UITitle, false)
				// update indicator current resource we expanded to
				v.expandIndicator.UpdatePath("Load Balancer", false)
			})
		}()

		// switch to browse page
		v.pages.SwitchToPage("browse")
	})

	v.startFocus = startScreen
	v.pages.AddPage("start", startScreen, true, true)
}

func (v *ViewController) createBrowseScreen(onSelect abstracts.OnTableSelectionFunc) tview.Primitive {

	// create grid
	grid := tview.NewGrid()
	grid.SetRows(6, 1, 0, 3).
		SetBorders(true)

	// header
	header := v.createHeaderView()
	// add header to grid
	headerInput := abstracts.NewDefaultGridItem(header)
	// header, headerGridRowIdx, 0, 1, 3, 0, 0, false
	headerInput.Metadata.Row = headerGridRowIdx
	grid.AddItem(headerInput.ConvertToAddItemInput())

	// next
	// next resource
	i am here
	nextIndicator := GetNextResourcesView()
	nextIndicator.SetNextAvailableResource(&NextResourcesInput{
		// Next: []ResourceInput{
		// 	{Type: awsUtils.ALBOrCLBType, UIName: "Load Balancer"},
		// 	{Type: awsUtils.LambdaType, UIName: "Lambda"},
		// 	{Type: awsUtils.EC2Type, UIName: "EC2"},
		// },
	})
	nextResourceView := nextIndicator.Get()
	nextInput := abstracts.NewDefaultGridItem(nextResourceView)
	nextInput.Metadata.Row = nextGridRowIdx
	grid.AddItem(nextInput.ConvertToAddItemInput())

	// table
	focusView := v.createFocusView(onSelect)
	v.focusView = focusView
	table := focusView.Get()

	// add table to grid
	tableInput := abstracts.NewDefaultGridItem(table)
	tableInput.Metadata.Row = focusViewGridIdx
	tableInput.Metadata.Focus = true
	grid.AddItem(tableInput.ConvertToAddItemInput())

	// footer

	footerWidget, expandIndicator := v.createFooterView("")

	v.expandIndicator = expandIndicator
	//here the problem, it never shows after render haha
	footer := footerWidget.GetRender()

	// add footer to grid
	footerInput := abstracts.NewDefaultGridItem(footer)
	footerInput.Metadata.Row = footerViewGridIdx
	grid.AddItem(footerInput.ConvertToAddItemInput())

	return grid
}

func (v *ViewController) createFocusView(onSelect abstracts.OnTableSelectionFunc) FocusView {
	desc := &abstracts.TableDescription{
		OnTableSelectionCallback: func(tsr *abstracts.TableSelectionResult) {
			onSelect(tsr)
		},
	}
	return NewFocusView(desc)
}

func (v *ViewController) createHeaderView() tview.Primitive {
	i := &HeaderInput{
		AppName: "TravelOps",
		AppFields: map[string]string{
			"AWS Profile": "default",
			"AWS Account": "12312312",
		},
	}

	header := NewHeaderView(i)
	return header
}

func (v *ViewController) createFooterView(initialResource string) (*IndicatorWidget, *ExpansionIndicatorManager) {
	footer := NewIndicatorView(initialResource)
	indicator := NewExpansionPathIndicator(initialResource, func(path []string, back bool) {
		p := ""

		if len(path) > 0 {
			p = path[len(path)-1]
		}
		footer.OnPathChange(p, back)
	})

	return footer, indicator
}

func (v *ViewController) createStartScreen(onSubmit func(out *StartPageOutput)) tview.Primitive {

	startInput := &StartPageInput{
		Title:       "Resource Starting Page",
		Description: "Choose a resource to start exploring your aws infra",
		AWSProfiles: []string{"default"},
		Options:     []*ResourceOption{{Type: awsUtils.R53RecordSetType, UITitle: "Route 53", UIDescription: "Start from DNS Address"}},
	}

	startView := NewStartForm(startInput, func(out *StartPageOutput) {
		onSubmit(out)
	})

	return startView
}

func (v *ViewController) Run() {
	if err := v.app.SetRoot(v.pages, true).SetFocus(v.startFocus).Run(); err != nil {
		panic(err)
	}
}

const pageCount = 5

func dummyPages() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	for page := 0; page < pageCount; page++ {
		func(page int) {
			pages.AddPage(fmt.Sprintf("page-%d", page),
				tview.NewModal().
					SetText(fmt.Sprintf("This is page %d. Choose where to go next.", page+1)).
					AddButtons([]string{"Next", "Quit"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						if buttonIndex == 0 {
							pages.SwitchToPage(fmt.Sprintf("page-%d", (page+1)%pageCount))
						} else {
							app.Stop()
						}
					}),
				false,
				page == 0)
		}(page)
	}
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
