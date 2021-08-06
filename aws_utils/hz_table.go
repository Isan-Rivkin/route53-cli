package aws_utils

import (
	"fmt"
)

const (
	HZNameCol         = "Hosted Zone"
	HZIdCol           = "Id"
	HZTotalRecordsCol = "Total records"
	HZPrivateCol      = "Private"
)

var supportedCols = []string{HZNameCol, HZIdCol, HZTotalRecordsCol, HZPrivateCol}

type HostedZoneR53ResultTableOutput struct {
	Outputs map[string]string
}

func (o *HostedZoneR53ResultTableOutput) GetHZTableCols() []string {
	return supportedCols
}

func (r *GetRecordAliasesResult) GetHostedZoneR53AsTableOutput() *HostedZoneR53ResultTableOutput {

	if r.HostedZone == nil {
		return nil
	}

	hz := r.HostedZone
	name := ""
	id := ""
	count := ""
	privateZone := ""

	if hz.Name != nil {
		name = *hz.Name
	}

	if hz.Id != nil {
		id = *hz.Id
	}

	if hz.ResourceRecordSetCount != nil {
		count = fmt.Sprint(*hz.ResourceRecordSetCount)
	}

	if hz.Config != nil && hz.Config.PrivateZone != nil {
		privateZone = fmt.Sprintf("%t", *hz.Config.PrivateZone)
	}

	hzTableoutput := map[string]string{
		HZNameCol:         name,
		HZIdCol:           id,
		HZTotalRecordsCol: count,
		HZPrivateCol:      privateZone,
	}

	return &HostedZoneR53ResultTableOutput{
		Outputs: hzTableoutput,
	}

}
