package aws_utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/ec2"
	elb "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/route53"
)

// resource <- query

// rtype <- identify(resource)

// info <- describe(rtype,resource)

// resourceContext <- contextInfo(info)

// display(resourceContext)

type ResourceType string

type Resource interface {
	ID() string
	Type() ResourceType
}

type AWSResourceIdentifier interface {
	InferFromRecordSet(r *route53.ResourceRecordSet) ([]ResourceType, error)
	InferTypeFromDNS(dns string) ([]ResourceType, error)
	InferRegionFromDNS(r *route53.ResourceRecordSet) string
	InferRegionFromResourceARN(arn string) string
}

type ResourceDescriber interface {
	Describe(rtype ResourceType, info interface{}) (interface{}, error)
}

type AWSResourceDescriber struct {
	session   *session.Session
	elbClient *elb.ELBV2
	acmClient *acm.ACM
	ec2Client *ec2.EC2
}

func NewAWSResourceDescriber(profile string) ResourceDescriber {
	return &AWSResourceDescriber{
		session: GetEnvSession(profile),
	}
}

func (desc *AWSResourceDescriber) Describe(rtype ResourceType, info interface{}) (interface{}, error) {
	var descErr error
	switch rtype {
	case ALBOrCLBType:
	case ELBType:
		input := info.(*LBDescriptionInput)
		return desc.describeLB(input)
	case TargetGroupType:
		input := info.(*TGDescriptionInput)
		return desc.describeTG(input)
	case ELBListenersType:
		input := info.(*LBListenersDescribeInput)
		return desc.describeELBListeners(input)
	case ACMType:
		input := info.(*ACMDescpInput)
		return desc.describeCertificate(input)
	case EC2Type:
		input := info.(*EC2InstanceDescInput)
		return desc.describeEC2Instances(input)
	default:
		descErr = fmt.Errorf("resource type not support for describe - %s", rtype)
	}

	return nil, descErr
}
