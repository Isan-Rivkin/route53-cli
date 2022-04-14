package aws_utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

type RoutingPolicy string

const (
	GeoLocationRP RoutingPolicy = "Geolocation"
	WeightedRP    RoutingPolicy = "Weighted"
	FailoverRP    RoutingPolicy = "Failover"
	OtherRP       RoutingPolicy = "-"
)

type RecordAccessor interface {
	GetHostedZoneID() string
	// "Record", "Type", "TTL", "Country", "Alias", "Resources"
	GetTTL() (int64, bool)
	// The record value in r53
	GetRecord() (string, bool)
	GetRegion() (string, bool)
	GetRecordType() (string, bool)
	// for GEO dns configurations
	GetCountryCode() (string, bool)
	// alias value is for A records, typical to ELB
	GetAliasValue() (string, bool)
	// resources are non alias values, typical for TXT, NS, MX, SOA records
	GetResources() ([]string, bool)
	// web URL to see the target Alias if exist
	GetWebURL() (string, bool)
	// DNS record feature types
	// Record has Routing policy and differentiator (i.e weighted-50 or geo-us, or failover-prinmary)
	// Failover, Simple, Geolocation, etc https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy.html
	GetRoutingPolicy() (RoutingPolicy, bool)
	GetRoutingPolicyDifferentiator() (string, bool)
	IsGeoDNS() bool
	IsWeightedDNS() bool
	IsFailoverRoutingPolicy() bool
	GetFailoverDifferentiator() (string, bool)
}

type RecordResultAccessor struct {
	raw *route53.ResourceRecordSet
	// from record set
	hostedZone *route53.HostedZone
	region     string
}

func FlatRecordResultToAccessors(results []*GetRecordAliasesResult) []RecordAccessor {
	var accessors []RecordAccessor
	for _, r := range results {
		accessors = append(accessors, r.GetRecordsAccessors()...)
	}
	return accessors
}

func (r *GetRecordAliasesResult) GetRecordsAccessors() []RecordAccessor {
	var accessors []RecordAccessor
	for _, recSet := range r.Records {
		ra := NewRecordResultAccessor(recSet, r.HostedZone, r.Region)
		accessors = append(accessors, ra)
	}
	return accessors
}

func NewRecordResultAccessor(raw *route53.ResourceRecordSet, hostedZone *route53.HostedZone, region string) RecordAccessor {
	return &RecordResultAccessor{
		raw:        raw,
		hostedZone: hostedZone,
		region:     region,
	}
}

func (a *RecordResultAccessor) GetHostedZoneID() string {
	return aws.StringValue(a.hostedZone.Id)
}

func (a *RecordResultAccessor) IsGeoDNS() bool {
	return a.raw.GeoLocation != nil
}

func (a *RecordResultAccessor) IsWeightedDNS() bool {
	return a.raw.Weight != nil
}

// RoutingPolicy Failover records
func (a *RecordResultAccessor) IsFailoverRoutingPolicy() bool {
	return aws.StringValue(a.raw.Failover) != ""
}

// RoutingPolicy Failover records: return PRIMARY / SECONDARY
func (a *RecordResultAccessor) GetFailoverDifferentiator() (string, bool) {
	return aws.StringValue(a.raw.Failover), a.IsFailoverRoutingPolicy()
}

func (a *RecordResultAccessor) GetRoutingPolicyDifferentiator() (string, bool) {

	if a.IsFailoverRoutingPolicy() {
		return a.GetFailoverDifferentiator()
	}

	if a.IsGeoDNS() {
		return aws.StringValue(a.raw.GeoLocation.SubdivisionCode), true
	}

	if a.IsWeightedDNS() {
		w := aws.Int64Value(a.raw.Weight)
		return fmt.Sprintf("%v", w), true
	}

	return "", true
}

func (a *RecordResultAccessor) GetRoutingPolicy() (RoutingPolicy, bool) {

	// differentiator contains primary/secondary etc values
	if a.IsFailoverRoutingPolicy() {
		return FailoverRP, true
	}
	// differentiator contains country
	if a.IsGeoDNS() {
		return GeoLocationRP, true
	}
	// differentiator contains weight number
	if a.IsWeightedDNS() {
		return WeightedRP, true
	}
	return OtherRP, false
}

func (a *RecordResultAccessor) GetTTL() (int64, bool) {
	return aws.Int64Value(a.raw.TTL), a.raw.TTL != nil
}

func (a *RecordResultAccessor) GetRegion() (string, bool) {
	var r string

	if a.raw.Region == nil {
		r = a.region
	} else {
		r = aws.StringValue(a.raw.Region)
	}

	return r, r != ""
}

func (a *RecordResultAccessor) GetRecord() (string, bool) {
	recordStr := aws.StringValue(a.raw.Name)
	if strings.HasPrefix(recordStr, WildCard) {
		recordStr = strings.Replace(recordStr, WildCard, "*", 1)
	}
	return recordStr, recordStr != ""
}

func (a *RecordResultAccessor) GetRecordType() (string, bool) {
	return aws.StringValue(a.raw.Type), a.raw.Type != nil
}

func (a *RecordResultAccessor) GetCountryCode() (string, bool) {
	countryCode := ""
	if a.raw.GeoLocation != nil {
		countryCode = aws.StringValue(a.raw.GeoLocation.SubdivisionCode)
	}

	return countryCode, countryCode != ""
}

func (a *RecordResultAccessor) GetAliasValue() (string, bool) {
	dnsName := ""
	if a.raw.AliasTarget != nil {
		dnsName = aws.StringValue(a.raw.AliasTarget.DNSName)
	}
	return dnsName, dnsName != ""
}
func (a *RecordResultAccessor) GetResources() ([]string, bool) {
	var resources []string
	for _, r := range a.raw.ResourceRecords {
		if r != nil {
			resources = append(resources, aws.StringValue(r.Value))
		}
	}
	return resources, len(resources) > 0
}

func (a *RecordResultAccessor) GetWebURL() (string, bool) {
	r, _ := a.GetRegion()
	a.raw.Region = aws.String(r)
	url, err := GenerateWebURL(a.raw)
	return url, err == nil
}
