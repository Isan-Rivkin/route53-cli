package v1

import (
	"fmt"
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rivo/tview"
)

type TableDescription struct {
	Headers []string
	Rows    []map[string]string
}

type ResourceTableDescriber func(resources []awsUtils.Resource) *TableDescription

var TableDescriber = map[awsUtils.ResourceType]ResourceTableDescriber{
	awsUtils.ALBOrCLBType: elbTableDescriber,
}

func (app *R53App) RenderFocusTable(rType awsUtils.ResourceType, resources []awsUtils.Resource) (*tview.Table, error) {
	tablePrompt := abstracts.NewTable()

	describerFn, exist := TableDescriber[rType]

	if !exist {
		return nil, fmt.Errorf("NoSuchResourceDescriber %s", rType)
	}

	description := describerFn(resources)

	header := description.Headers

	tablePrompt.AddHeaders(append([]string{"#"}, header...))

	currentRow := 0

	// populate columns
	for _, row := range description.Rows {

		cellRefValue := ""

		// add the # col for row number
		// populate cell
		tablePrompt.AddRow(currentRow+1, 0, fmt.Sprintf("%d", currentRow+1), cellRefValue)

		// add all the cols in the result
		for headerNum, headerName := range header {
			// populate cell
			tablePrompt.AddRow(currentRow+1, headerNum+1, row[headerName], cellRefValue)
		}
		currentRow++
	}

	return tablePrompt.Render(), nil
}

func elbTableDescriber(resources []awsUtils.Resource) *TableDescription {
	headers := []string{
		"Name",
		"Arn",
		"DNS",
		"Region",
		"State",
	}
	var rows []map[string]string

	for _, r := range resources {
		lbOutput := r.GetQueryOutput().(*awsUtils.LBDescriptionOutput)
		result, _ := r.TryGetValue(awsUtils.RegionAttr)
		region := result[0]

		for _, lb := range lbOutput.LoadBalancers {
			rowVals := map[string]string{
				"Name":   aws.StringValue(lb.LoadBalancerName),
				"Arn":    aws.StringValue(lb.LoadBalancerArn),
				"DNS":    aws.StringValue(lb.DNSName),
				"Region": region,
				"State":  aws.StringValue(lb.State.Code),
			}
			rows = append(rows, rowVals)
		}
	}
	return &TableDescription{
		Headers: headers,
		Rows:    rows,
	}

}
