package aws_utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
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
	InferRegionFromDNS(r *route53.ResourceRecordSet) string
}

type ResourceDescriber interface {
	Describe(rtype ResourceType, info interface{}) (interface{}, error)
}

type AWSResourceDescriber struct {
	session   *session.Session
	elbClient *elb.ELBV2
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
		output, err := desc.describeLB(input)
		return output, err
	case TargetGroupType:
		input := info.(*TGDescriptionInput)
		output, err := desc.describeTG(input)
		return output, err
	default:
		descErr = fmt.Errorf("resource type not support for describe - %s", rtype)
	}

	return nil, descErr
}
