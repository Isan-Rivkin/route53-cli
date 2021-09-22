package aws_utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

const (
	ErrNotRoutableAWS = "NotRoutableAWSResource"
	ErrGettingType    = "CantInferTypeErr"
)

const (
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
	ELBListenersType     ResourceType = "listener"
	EC2Type              ResourceType = "ec2"
	LambdaType           ResourceType = "lambda"
	RawIpType            ResourceType = "ip"
	ACMType              ResourceType = "acm"
)

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

type DefaultResourceIdentifier struct {
}

func NewDefaultResourceIdentifier() AWSResourceIdentifier {
	return &DefaultResourceIdentifier{}
}

func (ri *DefaultResourceIdentifier) InferFromRecordSet(r *route53.ResourceRecordSet) ([]ResourceType, error) {
	dnsType, routable := CheckRoutableAWSTarget(r)

	if !routable {
		return nil, fmt.Errorf(ErrNotRoutableAWS)
	}

	resourceType, found := dnsTargetsToTypes[dnsType]

	if !found {
		return nil, fmt.Errorf("%s - %s not supported for record %s", ErrGettingType, dnsType, *r.AliasTarget.DNSName)
	}

	return []ResourceType{resourceType}, nil
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
