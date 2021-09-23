package cliui

import (
	"fmt"
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"

	"github.com/rivo/tview"
)

func (app *R53App) OnR53TableSelection(selection *abstracts.TableSelectionResult) {

	app.EventsController <- &AppEvent{
		Type:         R53TableSelection,
		EventPayload: selection,
	}
}

func (app *R53App) SetR53RecordsQueryResult(result *awsUtils.GetRecordAliasesResult) {
	// get records gui table
	recordsTable, output := app.RenderR53RecordsTable(result, app.OnR53TableSelection)
	// get hosted zone text view
	hostedZoneLabeled := app.RenderHostedZoneTextView(output)
	hzTextView := hostedZoneLabeled.Render()
	app.initDefaultRootGrid(hzTextView, recordsTable)
}

func (app *R53App) RenderR53RecordsTable(result *awsUtils.GetRecordAliasesResult, onSelected abstracts.OnTableSelectionFunc) (*tview.Table, *awsUtils.R53ResultTableOutput) {
	tablePrompt := abstracts.NewTable()

	tablePrompt.AddSelectionCallBack(onSelected)

	// here web url is truncated awsUtils.WebURLCol since it's too big for the ui col
	// TODO:: adding another row will not help for url since it will expend it's width for the whole table
	// maybe create multiple tables with rows under it w/e
	header := []string{awsUtils.NumCol, awsUtils.RecordCol, awsUtils.TypeCol, awsUtils.TTLCol, awsUtils.CountryCol, awsUtils.AliasCol, awsUtils.ResourcesCol}

	output := result.GetR53AsTableOutput(&awsUtils.R53ResultTableInput{
		Columns: header,
	})

	nonEmptyHeaders := output.GetNonEmptyCols()

	tablePrompt.AddHeaders(append([]string{"#"}, nonEmptyHeaders...))

	currentRow := 0

	// check if alias column exist for reference value
	aliasHeaderExist := false
	for _, h := range nonEmptyHeaders {
		if h == awsUtils.AliasCol {
			aliasHeaderExist = true
			break
		}
	}
	// populate columns
	for _, row := range output.Outputs {

		cellRefValue := ""

		if aliasHeaderExist {
			cellRefValue = row[awsUtils.AliasCol]
		}

		// add the # col for row number
		// populate cell
		tablePrompt.AddRow(currentRow+1, 0, fmt.Sprintf("%d", currentRow+1), cellRefValue)

		// add all the cols in the result
		for headerNum, headerName := range nonEmptyHeaders {
			// populate cell
			tablePrompt.AddRow(currentRow+1, headerNum+1, row[headerName], cellRefValue)
		}

		// if web url exist create additional row below with the result
		// !!! this will not work because the column width of url will be the size of the web url so display will still be the same with additional empty row
		// if url, found := row[awsUtils.WebURLCol]; found {
		// 	tablePrompt.AddInfoRow(currentRow+2, 0, url)
		// 	currentRow++
		// }

		currentRow++
	}

	return tablePrompt.Render(), output
}

// SelectR53RecordFromList is return selected route53 record from the list by prompt question
func DEPRECATEDSelectR53RecordFromList(result *awsUtils.GetRecordAliasesResult) (string, error) {

	header := []string{awsUtils.NumCol, awsUtils.RecordCol, awsUtils.TypeCol, awsUtils.TTLCol, awsUtils.CountryCol, awsUtils.AliasCol, awsUtils.ResourcesCol}

	output := result.GetR53AsTableOutput(&awsUtils.R53ResultTableInput{
		Columns: header,
	})

	nonEmptyHeaders := output.GetNonEmptyCols()

	tablePrompt := NewTable()

	tablePrompt.AddHeaders(append([]string{"#"}, nonEmptyHeaders...))

	for i, row := range output.Outputs {

		tablePrompt.AddRow(i+1, 0, fmt.Sprintf("%d", i+1))

		for headerNum, headerName := range nonEmptyHeaders {
			tablePrompt.AddRow(i+1, headerNum+1, row[headerName])
		}

	}

	selectedIdx, err := tablePrompt.Render()
	if err != nil {
		return "", err
	}

	selectedRow := output.Outputs[selectedIdx]

	alias := selectedRow[awsUtils.AliasCol]

	return alias, nil
}
