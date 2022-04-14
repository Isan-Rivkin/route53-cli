package aws_utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/v6/table"
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
	accessors := r.GetRecordsAccessors()

	for _, rec := range accessors {
		countryCode, _ := rec.GetCountryCode()
		dnsName, _ := rec.GetAliasValue()
		ttl, _ := rec.GetTTL()
		resources, _ := rec.GetResources()
		resourcesRow := strings.Join(resources, "\n")
		recordStr, _ := rec.GetRecord()
		if opts != nil && opts.WebURL {
			url, _ := rec.GetWebURL()
			dnsName += fmt.Sprintf("\n\n%s\n", url)
		}
		recType, _ := rec.GetRecordType()
		t.AppendRow(table.Row{recordStr, recType, ttl, countryCode, dnsName, resourcesRow}, rowConfigAutoMerge)
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

	hzIdDisplay := aws.StringValue(r.HostedZone.Id)
	r53Url := GenerateRoute53HostedZoneWebURL(hzIdDisplay)
	if r53Url != "" {
		hzIdDisplay = fmt.Sprintf("%s\n\n%s\n", hzIdDisplay, r53Url)
	}
	t.AppendHeader(table.Row{"Hosted Zone", "Id", "Total records", "Private"}, rowConfigAutoMerge)
	t.AppendRow(table.Row{*r.HostedZone.Name, hzIdDisplay, *r.HostedZone.ResourceRecordSetCount, *r.HostedZone.Config.PrivateZone}, rowConfigAutoMerge)
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
