package aws_utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

func NewR53RecordSetDescriber(dns, region string) Resource {
	keys := map[ResourceKey][]string{}
	keys[DNSAttr] = []string{dns}
	keys[RegionAttr] = []string{region}
	return NewResourceDescriptionResult(nil, nil, nil, R53RecordSetType, keys, dns)
}

func NewR53RecordSetFromDNS(dns, region string) *route53.ResourceRecordSet {
	return &route53.ResourceRecordSet{
		Region: aws.String(region),
		AliasTarget: &route53.AliasTarget{
			DNSName: aws.String(dns),
		},
	}
}
