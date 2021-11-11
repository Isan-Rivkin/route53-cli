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
	UITitle     string
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

	form.AddInputField("Search Query", "", 100, nil, func(text string) {
		searchQuery = text
	})

	form.SetBorder(true).SetTitle("Enter Starting Point").SetTitleAlign(tview.AlignLeft)

	form.AddButton("Expand", func() {
		output := &StartPageOutput{AWSProfile: awsProfile}
		for _, i := range input.Options {
			if i.UITitle == startResource {
				output.Type = i.Type
				output.UITitle = i.UITitle
			}
		}
		output.SearchQuery = searchQuery
		onSelection(output)
	})

	form.AddButton("Quit", func() {

	})

	return form
}
