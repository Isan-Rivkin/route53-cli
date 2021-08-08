package cliui

import (
	"fmt"
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"

	"github.com/rivo/tview"
)

func (app *R53App) SetR53RecordsQueryResult(result *awsUtils.GetRecordAliasesResult) {
	// get records gui table
	recordsTable, output := app.RenderR53RecordsTable(result, nil)
	// get hosted zone text view
	hostedZoneLabeled := app.RenderHostedZoneTextView(output)
	hzTextView := hostedZoneLabeled.Render()
	app.initDefaultRootGrid(hzTextView, recordsTable)
}

func (app *R53App) RenderR53RecordsTable(result *awsUtils.GetRecordAliasesResult, onSelected abstracts.OnTableSelectionFunc) (*tview.Table, *awsUtils.R53ResultTableOutput) {
	tablePrompt := abstracts.NewTable()

	tablePrompt.AddSelectionCallBack(onSelected)

	header := []string{awsUtils.NumCol, awsUtils.RecordCol, awsUtils.TypeCol, awsUtils.TTLCol, awsUtils.CountryCol, awsUtils.AliasCol, awsUtils.ResourcesCol}

	output := result.GetR53AsTableOutput(&awsUtils.R53ResultTableInput{
		Columns: header,
	})

	nonEmptyHeaders := output.GetNonEmptyCols()

	tablePrompt.AddHeaders(append([]string{"#"}, nonEmptyHeaders...))

	currentRow := 0
	for _, row := range output.Outputs {

		// add the # col for row number
		tablePrompt.AddRow(currentRow+1, 0, fmt.Sprintf("%d", currentRow+1))

		// add all the cols in the result
		for headerNum, headerName := range nonEmptyHeaders {
			tablePrompt.AddRow(currentRow+1, headerNum+1, row[headerName])
		}

		// if web url exist create additional row below with the result
		// if url, found := row[awsUtils.WebURLCol]; found {
		// 	tablePrompt.AddInfoRow(currentRow+2, 0, url)
		// 	currentRow++
		// }

		currentRow++
	}

	return tablePrompt.Render(), output
}

// SelectR53RecordFromList is return selected route53 record from the list by prompt question
func SelectR53RecordFromList(result *awsUtils.GetRecordAliasesResult) (string, error) {

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
