package aws_utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

const (
	ErrNotRoutableAWS = "NotRoutableAWSResource"
	ErrGettingType    = "CantInferTypeErr"
)

const (
	R53RecordSetType     ResourceType = "r53.record"
	ALBType              ResourceType = "lb.alb"
	CLBType              ResourceType = "lb.clb"
	ALBOrCLBType         ResourceType = "lb.alb_clb"
	NLBType              ResourceType = "lb.nlb"
	ELBType              ResourceType = "lb.elb"
	S3Type               ResourceType = "s3"
	VPCEndpointType      ResourceType = "vpce"
	CloudFrontType       ResourceType = "cloudfront"
	ElasticBeansTalkType ResourceType = "elasticbeanstalk"
	AcceleratorApiType   ResourceType = "awsglobalaccelerator"
	TargetGroupType      ResourceType = "targetgroup"
	HealthTargetsTGType  ResourceType = "tg.targets.health"
	ELBListenersType     ResourceType = "listener"
	EC2Type              ResourceType = "ec2"
	LambdaType           ResourceType = "lambda"
	RawIpType            ResourceType = "ip"
	ACMType              ResourceType = "acm"
	SecurityGroupType    ResourceType = "sg"
)

// reachable resources from source resource, i.e a graph of potential connectivity
var AdjacentResources = map[ResourceType][]ResourceType{
	R53RecordSetType:    {ALBType, CLBType, ELBType, ALBOrCLBType},
	ALBType:             {TargetGroupType, ELBListenersType},
	CLBType:             {TargetGroupType, ELBListenersType},
	NLBType:             {TargetGroupType, ELBListenersType},
	ELBType:             {TargetGroupType, ELBListenersType},
	ALBOrCLBType:        {TargetGroupType, ELBListenersType},
	TargetGroupType:     {HealthTargetsTGType},
	HealthTargetsTGType: {EC2Type, LambdaType, RawIpType},
	ELBListenersType:    {ACMType},
	ACMType:             {},
	EC2Type:             {},
	LambdaType:          {},
	RawIpType:           {},
	// not implemented yet
	SecurityGroupType:    {},
	S3Type:               {},
	VPCEndpointType:      {},
	CloudFrontType:       {},
	ElasticBeansTalkType: {},
	AcceleratorApiType:   {},
}

var dnsTargetsToTypes = map[string]ResourceType{
	ALBDns:              ALBOrCLBType,
	NLBDns:              NLBType,
	ELBDns:              ELBType,
	S3Suffx:             S3Type,
	VPCEndpointDns:      VPCEndpointType,
	CloudFrontDns:       CloudFrontType,
	ElasticBeansTalkDns: ElasticBeansTalkType,
	AcceleratorApiDns:   AcceleratorApiDns,
}

// check if a resource can be reached from another resource
func IsReachableFrom(from ResourceType, to ResourceType) bool {
	links, found := AdjacentResources[from]

	if !found {
		return false
	}

	for _, l := range links {
		if l == to {
			return true
		}
	}
	return false
}

type AWSResourceIdentifier interface {
	InferFromRecordSet(r *route53.ResourceRecordSet) ([]ResourceType, error)
	InferTypeFromDNS(dns string) ([]ResourceType, error)
	InferRegionFromDNS(r *route53.ResourceRecordSet) string
	InferRegionFromResourceARN(arn string) string
}

type DefaultResourceIdentifier struct {
}

func NewDefaultResourceIdentifier() AWSResourceIdentifier {
	return &DefaultResourceIdentifier{}
}

func (ri *DefaultResourceIdentifier) InferTypeFromDNS(dns string) ([]ResourceType, error) {
	dnsType, routable := checkRoutableAWSTargetDNS(dns)

	if !routable {
		return nil, fmt.Errorf(ErrNotRoutableAWS)
	}

	resourceType, found := dnsTargetsToTypes[dnsType]

	if !found {
		return nil, fmt.Errorf("%s - %s not supported for record %s", ErrGettingType, dnsType, dns)
	}

	return []ResourceType{resourceType}, nil
}
func (ri *DefaultResourceIdentifier) InferFromRecordSet(r *route53.ResourceRecordSet) ([]ResourceType, error) {
	if r.AliasTarget == nil || r.AliasTarget.DNSName == nil {
		return nil, errors.New("ErrNoDNSTarget")
	}

	dns := *r.AliasTarget.DNSName

	return ri.InferTypeFromDNS(dns)
}

func (ri *DefaultResourceIdentifier) InferRegionFromDNS(r *route53.ResourceRecordSet) string {
	dnsType, routable := CheckRoutableAWSTarget(r)

	if !routable {
		return ""
	}

	return GetRegionFromLBDNSName(dnsType, r)
}

// arn is always arn:partition:service:region:account-id:resource-type:resource-id
func (ri *DefaultResourceIdentifier) InferRegionFromResourceARN(arn string) string {

	if !strings.HasPrefix(arn, "arn:") {
		return ""
	}

	arnChunks := strings.Split(arn, ":")
	if len(arnChunks) < 4 {
		return ""
	}

	return arnChunks[3]
}
