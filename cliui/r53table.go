package cliui

import (
	"fmt"
	awsUtils "r53/aws_utils"
)

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
