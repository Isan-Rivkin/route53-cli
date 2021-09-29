package cliui

import (
	awsUtils "r53/aws_utils"

	"github.com/rivo/tview"
)

type ResourceOption struct {
	Type          awsUtils.ResourceType
	UITitle       string
	UIDescription string
}

type StartPageInput struct {
	Title       string
	Description string
	Options     []*ResourceOption
	AWSProfiles []string
}
type StartPageOutput struct {
	AWSProfile  string
	SearchQuery string
	Type        awsUtils.ResourceType
}

func NewStartForm(input *StartPageInput, onSelection func(out *StartPageOutput)) *tview.Form {

	var (
		dropdown                               []string
		awsProfile, startResource, searchQuery string
	)

	for _, i := range input.Options {
		dropdown = append(dropdown, i.UITitle)
	}

	form := tview.NewForm()

	form.AddDropDown("AWS Profile", input.AWSProfiles, 0, func(option string, optionIndex int) {
		awsProfile = option
	})

	awsProfile = "default"

	form.AddDropDown("Starting Resource", dropdown, 0, func(option string, optionIndex int) {
		startResource = option
	})

	form.AddInputField("Search Query", "", 200, nil, func(text string) {
		searchQuery = text
	})

	form.SetBorder(true).SetTitle("Enter Starting Point").SetTitleAlign(tview.AlignLeft)

	form.AddButton("Expand", func() {
		output := &StartPageOutput{AWSProfile: awsProfile}
		for _, i := range input.Options {
			if i.UITitle == startResource {
				output.Type = i.Type
			}
		}
		output.SearchQuery = searchQuery
		onSelection(output)
	})

	form.AddButton("Quit", func() {

	})

	return form
}

func NewStartView(input *StartPageInput, onSelection func(out *StartPageOutput)) tview.Primitive {
	// Returns a new primitive which puts the provided primitive in the center and
	// sets its size to the given width and height.
	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, false).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false)
	}

	form := NewStartForm(input, onSelection)

	return form
	grid := tview.NewGrid().
		SetRows(30).
		SetColumns(30).
		SetBorders(true).
		AddItem(form, 0, 0, 1, 1, 0, 0, true)

	return grid
	return modal(form, 40, 10)
	pages := tview.NewPages().
		AddPage("start", modal(form, 40, 10), true, true)

	return pages
}
