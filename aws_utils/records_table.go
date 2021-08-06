package aws_utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
)

const (
	RecordCol    = "Record"
	TypeCol      = "Type"
	NumCol       = "#"
	TTLCol       = "TTL"
	CountryCol   = "Country"
	AliasCol     = "Alias"
	ResourcesCol = "Resources"
	WebURLCol    = "URL"
)

type R53ResultTableOutput struct {
	HzOutput     *HostedZoneR53ResultTableOutput
	Outputs      []map[string]string
	NonEmptyCols map[string]bool
}

func (o *R53ResultTableOutput) GetNonEmptyCols() []string {
	var cols []string
	for colName, _ := range o.NonEmptyCols {
		cols = append(cols, colName)
	}
	return cols
}

func (o *R53ResultTableOutput) IsEmptyCol(colName string) bool {
	_, found := o.NonEmptyCols[colName]
	return !found
}

type R53ResultTableInput struct {
	PrintOpts *PrintOptions
	Columns   []string
}

func (r *GetRecordAliasesResult) getSupportedParsers() map[string]func(*route53.ResourceRecordSet) (string, bool) {

	m := map[string]func(*route53.ResourceRecordSet) (string, bool){
		RecordCol:    r.recordCol,
		TypeCol:      r.recordType,
		TTLCol:       r.recordTTL,
		CountryCol:   r.recordCountry,
		AliasCol:     r.recordAlias,
		ResourcesCol: r.recordResources,
		WebURLCol:    r.recordURL,
	}

	return m
}

func (r *GetRecordAliasesResult) GetR53AsTableOutput(input *R53ResultTableInput) *R53ResultTableOutput {

	// parse r53 recourds results into a dynamic table

	output := &R53ResultTableOutput{}

	nonEmptyCols := map[string]bool{}

	supportedCols := r.getSupportedParsers()

	for _, recSet := range r.Records {

		outputRow := map[string]string{}

		if recSet.Region == nil {
			recSet.Region = &r.Region
		}

		for _, col := range input.Columns {
			if parser, supported := supportedCols[col]; supported {
				if val, found := parser(recSet); found {
					nonEmptyCols[col] = true
					outputRow[col] = val
				}

			}
		}

		output.Outputs = append(output.Outputs, outputRow)
	}

	output.NonEmptyCols = nonEmptyCols

	// parse hostedZone table

	output.HzOutput = r.GetHostedZoneR53AsTableOutput()

	return output
}

func (r *GetRecordAliasesResult) recordCol(recSet *route53.ResourceRecordSet) (string, bool) {
	recordStr := ""
	found := false
	if recSet.Name != nil {
		found = true
		recordStr = *recSet.Name
		if strings.HasPrefix(recordStr, WildCard) {
			recordStr = strings.Replace(recordStr, WildCard, "*", 1)
		}
	}
	return recordStr, found
}
func (r *GetRecordAliasesResult) recordType(recSet *route53.ResourceRecordSet) (string, bool) {
	return *recSet.Type, recSet.Type != nil
}

func (r *GetRecordAliasesResult) recordTTL(recSet *route53.ResourceRecordSet) (string, bool) {
	ttl := int64(0)
	if recSet.TTL != nil {
		ttl = *recSet.TTL
	}
	return fmt.Sprint(ttl), ttl != int64(0)
}

func (r *GetRecordAliasesResult) recordCountry(recSet *route53.ResourceRecordSet) (string, bool) {
	countryCode := ""
	if recSet.GeoLocation != nil && recSet.GeoLocation.SubdivisionCode != nil {
		countryCode = *recSet.GeoLocation.SubdivisionCode
	}
	return countryCode, countryCode != ""
}

func (r *GetRecordAliasesResult) recordAlias(recSet *route53.ResourceRecordSet) (string, bool) {
	dnsName := ""
	if recSet.AliasTarget != nil && recSet.AliasTarget.DNSName != nil {
		dnsName = *recSet.AliasTarget.DNSName
	}
	return dnsName, dnsName != ""
}

func (r *GetRecordAliasesResult) recordURL(recSet *route53.ResourceRecordSet) (string, bool) {

	url, err := GenerateWebURL(recSet)

	if err != nil {
		log.WithField("record", *recSet.Name).WithError(err).Debug("failed getting web url for record")
	}

	return url, err == nil && url != ""
}

func (r *GetRecordAliasesResult) recordResources(recSet *route53.ResourceRecordSet) (string, bool) {
	resourcesRow := ""
	for _, resources := range recSet.ResourceRecords {
		if resources != nil {
			resourcesRow += *resources.Value + "\n"
		}
	}
	return resourcesRow, resourcesRow != ""
}
