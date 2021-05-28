package aws_utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
)

type PrintOptions struct {
	WebURL bool
}

func (r *GetRecordAliasesResult) PrintTable(opts *PrintOptions) {
	r.printHostedzoneTable(opts)
	r.printRecordsTable(opts)
}

func (r *GetRecordAliasesResult) printRecordsTable(opts *PrintOptions) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Record", "Type", "TTL", "Country", "Alias", "Resources"}, rowConfigAutoMerge)
	for _, recSet := range r.Records {

		if recSet.Region == nil {
			recSet.Region = &r.Region
		}

		countryCode := ""
		if recSet.GeoLocation != nil && recSet.GeoLocation.SubdivisionCode != nil {
			countryCode = *recSet.GeoLocation.SubdivisionCode
		}
		dnsName := ""
		if recSet.AliasTarget != nil && recSet.AliasTarget.DNSName != nil {
			dnsName = *recSet.AliasTarget.DNSName
		}
		ttl := int64(0)
		if recSet.TTL != nil {
			ttl = *recSet.TTL
		}
		resourcesRow := ""
		for _, resources := range recSet.ResourceRecords {
			if resources != nil {
				resourcesRow += *resources.Value + "\n"
			}
		}
		recordStr := ""
		if recSet.Name != nil {
			recordStr = *recSet.Name
			if strings.HasPrefix(recordStr, WildCard) {
				recordStr = strings.Replace(recordStr, WildCard, "*", 1)
			}
		}
		// check if web url should be added
		if opts != nil && opts.WebURL {

			url, err := GenerateWebURL(recSet)
			if err == nil {
				dnsName += fmt.Sprintf("\n\n%s\n", url)
			} else {
				log.WithField("record", *recSet.Name).WithError(err).Debug("failed getting web url for record")
			}
		}
		t.AppendRow(table.Row{recordStr, *recSet.Type, ttl, countryCode, dnsName, resourcesRow}, rowConfigAutoMerge)
		t.AppendSeparator()
	}
	t.SetAutoIndex(true)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: false},
		{Number: 3, AutoMerge: false},
		{Number: 4, AutoMerge: false},
		{Number: 5, AutoMerge: false},
	})
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	t.SetOutputMirror(os.Stdout)
	t.Render()
}
func (r *GetRecordAliasesResult) printHostedzoneTable(opts *PrintOptions) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	t := table.NewWriter()

	t.AppendHeader(table.Row{"Hosted Zone", "Id", "Total records", "Private"}, rowConfigAutoMerge)
	t.AppendRow(table.Row{*r.HostedZone.Name, *r.HostedZone.Id, *r.HostedZone.ResourceRecordSetCount, *r.HostedZone.Config.PrivateZone}, rowConfigAutoMerge)
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
		{Number: 2, AutoMerge: true},
		{Number: 3, AutoMerge: true},
		{Number: 4, AutoMerge: true},
	})
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	t.SetOutputMirror(os.Stdout)
	t.Render()
}
