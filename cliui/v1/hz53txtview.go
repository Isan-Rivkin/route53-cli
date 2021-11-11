package v1

import (
	"fmt"
	awsUtils "r53/aws_utils"
	"r53/cliui/abstracts"
)

func (app *R53App) RenderHostedZoneTextView(resultTableOut *awsUtils.R53ResultTableOutput) *abstracts.LabeledText {

	hz := resultTableOut.HzOutput.Outputs

	// name
	header := abstracts.NewLabeledText(&abstracts.LabeledTextViewInput{
		ViewTitle: "Hosted Zone Info",
		Rows: []*abstracts.LabeledTextRowInput{
			{
				LabelTxt:   awsUtils.HZNameCol,
				LabelColor: abstracts.OrangeTxtColor,
				ValColor:   abstracts.WhiteTxtColor,
				ValTxt:     hz[awsUtils.HZNameCol],
			},
			{
				LabelTxt:   awsUtils.HZIdCol,
				LabelColor: abstracts.OrangeTxtColor,
				ValTxt:     hz[awsUtils.HZIdCol],
				ValColor:   abstracts.WhiteTxtColor,
			},
			{
				LabelTxt:   awsUtils.HZTotalRecordsCol,
				LabelColor: abstracts.OrangeTxtColor,
				ValTxt:     hz[awsUtils.HZTotalRecordsCol],
				ValColor:   abstracts.WhiteTxtColor,
			},
			{
				LabelTxt:   awsUtils.HZPrivateCol,
				LabelColor: abstracts.OrangeTxtColor,
				ValTxt:     hz[awsUtils.HZPrivateCol],
				ValColor:   abstracts.WhiteTxtColor,
			},
			{
				LabelTxt:   "NS Verified:",
				LabelColor: abstracts.OrangeTxtColor,
				ValColor:   abstracts.WhiteTxtColor,
				ValTxt:     fmt.Sprint(resultTableOut.HzOutput.VerifiedNameservers),
			},
		},
	})

	return header

}
