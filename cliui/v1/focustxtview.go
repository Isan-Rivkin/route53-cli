package v1

import (
	"fmt"
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"

	"github.com/aws/aws-sdk-go/aws"
)

type ResourceTxtViewDescription struct {
	Title           string
	LabelColor      string
	ValueColor      string
	Labels          []string
	LabelsToRowItem []map[string]string
}

type ResourceTxtViewDescriber func(resources []awsUtils.Resource) (*ResourceTxtViewDescription, error)

var TxtViewDescriber = map[awsUtils.ResourceType]ResourceTxtViewDescriber{
	awsUtils.ALBOrCLBType: elbTxtViewDescriber,
}

func (app *R53App) RenderFocusTxtViewItems(rType awsUtils.ResourceType, resources []awsUtils.Resource) (*abstracts.LabeledText, error) {
	describerFn, exist := TxtViewDescriber[rType]

	if !exist {
		return nil, fmt.Errorf("NoSuchResourceDescriber %s", rType)
	}

	description, err := describerFn(resources)

	if err != nil {
		return nil, err
	}

	var view *abstracts.LabeledText

	viewInput := &abstracts.LabeledTextViewInput{
		ViewTitle: description.Title,
		Rows:      []*abstracts.LabeledTextRowInput{},
	}

	for _, row := range description.LabelsToRowItem {

		for _, label := range description.Labels {
			value := row[label]

			viewInput.Rows = append(viewInput.Rows, &abstracts.LabeledTextRowInput{
				LabelTxt:   label,
				LabelColor: description.LabelColor,
				ValColor:   description.ValueColor,
				ValTxt:     value,
			})

		}

	}

	view = abstracts.NewLabeledText(viewInput)
	return view, nil
}

func elbTxtViewDescriber(resources []awsUtils.Resource) (*ResourceTxtViewDescription, error) {
	labels := []string{
		"Name",
		"URL",
		"Arn",
		"DNS",
		"Region",
		"State",
	}
	var rows []map[string]string

	for _, r := range resources {
		lbOutput := r.GetQueryOutput().(*awsUtils.LBDescriptionOutput)
		result, _ := r.TryGetValue(awsUtils.RegionAttr)
		webUrls, _ := r.TryGetValue(awsUtils.WebURLAttr)

		region := result[0]
		for idx, lb := range lbOutput.LoadBalancers {
			dns := aws.StringValue(lb.DNSName)
			rowVals := map[string]string{
				"Name":   aws.StringValue(lb.LoadBalancerName),
				"URL":    webUrls[idx],
				"Arn":    aws.StringValue(lb.LoadBalancerArn),
				"DNS":    dns,
				"Region": region,
				"State":  aws.StringValue(lb.State.Code),
			}
			rows = append(rows, rowVals)
		}
	}

	return &ResourceTxtViewDescription{
		Title:           "Load Balancers",
		LabelColor:      abstracts.OrangeTxtColor,
		ValueColor:      abstracts.WhiteTxtColor,
		Labels:          labels,
		LabelsToRowItem: rows,
	}, nil
}
