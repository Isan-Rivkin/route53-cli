package sdk

import (
	"encoding/json"

	awsu "github.com/isan-rivkin/route53-cli/aws_utils"

	"github.com/aws/aws-sdk-go/aws"
)

type RoutingOutput struct {
	Policy         awsu.RoutingPolicy `json:"policy,omitempty"`
	Differentiator string             `json:"differentiator,omitempty"`
}

type RecordOutput struct {
	HzID      string        `json:"hostedZoneId"`
	Record    string        `json:"record"`
	Type      string        `json:"type"`
	Alias     string        `json:"alias"`
	WebURL    string        `json:"url"`
	Routing   RoutingOutput `json:"routing,omitempty"`
	Resources []string      `json:"resources,omitempty"`
}

type HostedZoneOutput struct {
	Name         string `json:"name"`
	Id           string `json:"id"`
	WebUrl       string `json:"url"`
	TotalRecords int64  `json:"totalRecords"`
	Private      bool   `json:"private"`
}

func NewHostedZoneOutput(a *awsu.GetRecordAliasesResult) HostedZoneOutput {
	id := aws.StringValue(a.HostedZone.Id)
	o := HostedZoneOutput{
		Name:         aws.StringValue(a.HostedZone.Name),
		Id:           id,
		WebUrl:       awsu.GenerateRoute53HostedZoneWebURL(id),
		TotalRecords: aws.Int64Value(a.HostedZone.ResourceRecordSetCount),
		Private:      aws.BoolValue(a.HostedZone.Config.PrivateZone),
	}
	return o
}

type ResultOutput struct {
	Records     []RecordOutput              `json:"records"`
	HostedZones map[string]HostedZoneOutput `json:"hostedZones"`
}

func NewRecordOutput(a awsu.RecordAccessor) RecordOutput {
	o := RecordOutput{
		HzID: a.GetHostedZoneID(),
	}
	rec, _ := a.GetRecord()
	typ, _ := a.GetRecordType()
	alias, _ := a.GetAliasValue()
	webUrl, _ := a.GetWebURL()

	resources, _ := a.GetResources()

	o.Record = rec
	o.Type = typ
	o.Alias = alias
	o.WebURL = webUrl

	o.Resources = resources

	rp, hasRp := a.GetRoutingPolicy()

	if hasRp {
		rpd, _ := a.GetRoutingPolicyDifferentiator()
		o.Routing = RoutingOutput{
			Policy:         rp,
			Differentiator: rpd,
		}
	}
	return o
}
func ToSimpleOutput(resultSet []*awsu.GetRecordAliasesResult) *ResultOutput {
	var records []RecordOutput
	accessors := awsu.FlatRecordResultToAccessors(resultSet)
	for _, a := range accessors {
		records = append(records, NewRecordOutput(a))
	}

	// add hosted zone info
	seenHz := map[string]HostedZoneOutput{}
	for _, rs := range resultSet {
		hz := NewHostedZoneOutput(rs)
		// make sure there is no duplicates
		if _, seen := seenHz[hz.Id]; !seen {
			seenHz[hz.Id] = hz
		}

	}

	return &ResultOutput{Records: records, HostedZones: seenHz}
}

func ToJSONOutput(resultSet []*awsu.GetRecordAliasesResult, pretty bool) ([]byte, error) {
	result := ToSimpleOutput(resultSet)
	if pretty {
		return json.MarshalIndent(result, "", "    ")
	} else {
		return json.Marshal(result)
	}
}
