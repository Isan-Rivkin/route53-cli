package cliui

import (
	"fmt"
	awsUtils "r53/aws_utils"

	"github.com/rivo/tview"
)

type View interface {
	Run()
}

type ViewController struct {
}

func NewViewController() View {
	dummyInit()
	return &ViewController{}
}

func dummyInit() {
	startInput := &StartPageInput{
		Title:       "Resource Starting Page",
		Description: "Choose a resource to start exploring your aws infra",
		AWSProfiles: []string{"default"},
		Options:     []*ResourceOption{{Type: awsUtils.R53RecordSetType, UITitle: "Route 53", UIDescription: "Start from DNS Address"}},
	}
	startView := NewStartView(startInput, func(out *StartPageOutput) {})
	app := tview.NewApplication()
	pages := tview.NewPages()

	pages.AddPage("Entry", startView, true, true)

	if err := app.SetRoot(pages, true).SetFocus(startView).Run(); err != nil {
		panic(err)
	}
}
func (v *ViewController) Run() {
	//.	dummyPages()
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
