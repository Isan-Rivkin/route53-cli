package abstracts

import (
	"fmt"

	"github.com/rivo/tview"
)

const (
	OrangeTxtColor = "orange"
	WhiteTxtColor  = "white"
	RedTxtColor    = "red"
	GreenTxtColor  = "green"
	YellowTxtColor = "yellow"
	PurpleTxtColor = "purple"
)

type LabeledTextRowInput struct {
	LabelTxt   string
	LabelColor string
	ValTxt     string
	ValColor   string
}

type LabeledTextViewInput struct {
	ViewTitle string
	Rows      []*LabeledTextRowInput
}

type LabeledText struct {
	text  *tview.TextView
	input *LabeledTextViewInput
}

func NewLabeledText(input *LabeledTextViewInput) *LabeledText {
	t := tview.NewTextView()

	return &LabeledText{
		text:  t,
		input: input,
	}
}

func (t *LabeledText) Render() *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)

	txt := ""

	textView.SetMaxLines(len(t.input.Rows) + 1)

	for _, row := range t.input.Rows {

		lineTxt := fmt.Sprintf(`[%s]%s [%s]%s`, row.LabelColor, row.LabelTxt, row.ValColor, row.ValTxt)
		txt += fmt.Sprintf("%s \n", lineTxt)

	}

	textView.SetText(txt)
	textView.SetTitle(t.input.ViewTitle)
	textView.SetBorder(true)

	return textView
}
