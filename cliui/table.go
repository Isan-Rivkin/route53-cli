package cliui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
)

var (
	ErrCancelSelection = errors.New("selection canceled")
)

type TablePrompt struct {
	app   *tview.Application
	table *tview.Table
}

// NewTable create new table view
func NewTable() *TablePrompt {

	return &TablePrompt{
		app:   tview.NewApplication().EnableMouse(true),
		table: tview.NewTable().SetBorders(false),
	}
}

// AddHeaders add headers to the grid
func (tp *TablePrompt) AddHeaders(headers []string) {

	for i, headerText := range headers {
		tp.table.SetCell(0, i, tview.NewTableCell(strings.ToUpper(headerText)).
			SetSelectable(false).
			SetTextColor(tcell.ColorWhite).
			SetTransparency(true).
			SetAttributes(tcell.AttrBold).SetAlign(tview.AlignLeft))
	}

}

// AddRow adding new row to the grid
func (tp *TablePrompt) AddRow(row, column int, text string) {
	tp.table.SetCell(row, column, tview.NewTableCell(fmt.Sprintf("%s    ", text)).SetTextColor(tcell.ColorDarkGray).SetAlign(tview.AlignLeft))
}

// Render shows the table the the screen
func (tp *TablePrompt) Render() (int, error) {

	rowSelected := -1
	tp.table.SetSelectable(true, false)
	tp.table.SetSelectedFunc(func(row int, column int) {
		tp.table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		tp.table.SetSelectable(true, false)
		rowSelected = row
		tp.app.Stop()

	})

	appGrid := tview.NewGrid().SetBorders(true).AddItem(tp.table, 0, 0, 1, 3, 0, 0, true)

	if err := tp.app.SetRoot(appGrid, true).Run(); err != nil {
		log.WithError(err).Error("could not render table")
		return rowSelected, errors.New("could not render table")
	}

	if rowSelected == -1 {
		return rowSelected, ErrCancelSelection
	}

	return rowSelected, nil

}

// SelectInstanceFromList is return selected instance from the list by prompt question
func SelectInstanceFromList(instances []string) (string, error) {

	if len(instances) == 1 {
		return instances[0], nil
	}
	tablePrompt := NewTable()

	header := []string{
		"#",
		"Instance ID",
		"IP",
		"Role",
		"Environment",
		"Instance Type",
		"Platform",
		"Region",
	}

	tablePrompt.AddHeaders(header)

	for i, row := range instances {
		tablePrompt.AddRow(i+1, 0, fmt.Sprintf("%d", i+1))
		tablePrompt.AddRow(i+1, 1, row)
		tablePrompt.AddRow(i+1, 2, row)
		tablePrompt.AddRow(i+1, 3, row)
		tablePrompt.AddRow(i+1, 4, row)
		tablePrompt.AddRow(i+1, 5, row)
		tablePrompt.AddRow(i+1, 6, row)
		tablePrompt.AddRow(i+1, 7, row)

	}

	selectedInstance, err := tablePrompt.Render()
	if err != nil {
		return "", err
	}

	return instances[selectedInstance-1], nil
}
