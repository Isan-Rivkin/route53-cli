package cliui

import (
	"fmt"
	"r53/cliui/abstracts"
	"strings"

	"github.com/rivo/tview"
)

type KeyboardShortcut struct {
	Text string
	Keys []string
}

type HeaderInput struct {
	AppName   string
	AppFields map[string]string
	Shortcuts []KeyboardShortcut
}

func NewHeaderTxt(i *HeaderInput) tview.Primitive {
	var view *abstracts.LabeledText

	viewInput := &abstracts.LabeledTextViewInput{
		ViewTitle: i.AppName,
		Rows:      []*abstracts.LabeledTextRowInput{},
	}

	// meta info fields of the app such as aws profile and aws account
	for field, val := range i.AppFields {
		r := &abstracts.LabeledTextRowInput{
			LabelTxt:   field,
			LabelColor: abstracts.OrangeTxtColor,
			ValTxt:     val,
			ValColor:   abstracts.WhiteTxtColor,
		}
		viewInput.Rows = append(viewInput.Rows, r)
	}

	for _, sc := range i.Shortcuts {
		r := &abstracts.LabeledTextRowInput{
			LabelTxt:   sc.Text,
			LabelColor: abstracts.PurpleTxtColor,
			ValTxt:     fmt.Sprintf("( %s )", strings.Join(sc.Keys, "+")),
			ValColor:   abstracts.YellowTxtColor,
		}
		viewInput.Rows = append(viewInput.Rows, r)
	}

	view = abstracts.NewLabeledText(viewInput)

	return view.Render()
}

func NewHeaderView(i *HeaderInput) tview.Primitive {
	return NewHeaderTxt(i)
}
