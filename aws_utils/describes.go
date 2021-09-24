package aws_utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/ec2"
	elb "github.com/aws/aws-sdk-go/service/elbv2"
)

type ResourceDescriber interface {
	Describe(rtype ResourceType, info interface{}) (interface{}, error)
	DescribeResource(rtype ResourceType, info interface{}) (Resource, error)
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

// wrapper for Describe that returns common identifier wrapper struct for the result
func (desc *AWSResourceDescriber) DescribeResource(rtype ResourceType, info interface{}) (Resource, error) {
	return desc.describe(rtype, info)
}

// Describe some aws resource based on the type and info input as some concrete input
func (desc *AWSResourceDescriber) Describe(rtype ResourceType, info interface{}) (interface{}, error) {
	description, err := desc.describe(rtype, info)

	if err != nil {
		return nil, err
	}

	return description.GetQueryOutput(), nil
}

// does what Describe says
func (desc *AWSResourceDescriber) describe(rtype ResourceType, info interface{}) (Resource, error) {

	var (
		descErr error
		keys    map[ResourceKey][]string
		output  interface{}
		outId   string
	)

	switch rtype {
	case ALBOrCLBType:
	case ELBType:
		input := info.(*LBDescriptionInput)
		res, err := desc.describeLB(input)
		keys = res.GetKeys()
		outId = res.GetOutputID()
		descErr = err
		output = res
	case TargetGroupType:
		input := info.(*TGDescriptionInput)
		res, err := desc.describeTG(input)
		keys = res.GetKeys()
		outId = res.GetOutputID()
		output = res
		descErr = err
	case HealthTargetsTGType:
		input := info.(*TargetHealthDescTGInput)
		res, err := desc.describeTargetsHealth(input)
		keys = res.GetKeys()
		outId = res.GetOutputID()
		output = res
		descErr = err
	case ELBListenersType:
		input := info.(*LBListenersDescribeInput)
		res, err := desc.describeELBListeners(input)
		keys = res.GetKeys()
		outId = res.GetOutputID()
		output = res
		descErr = err
	case ACMType:
		input := info.(*ACMDescpInput)
		res, err := desc.describeCertificate(input)
		keys = res.GetKeys()
		outId = res.GetOutputID()
		output = res
		descErr = err
	case EC2Type:
		input := info.(*EC2InstanceDescInput)
		res, err := desc.describeEC2Instances(input)
		keys = res.GetKeys()
		outId = res.GetOutputID()
		output = res
		descErr = err
	default:
		descErr = fmt.Errorf("resource type not support for describe - %s", rtype)
		return nil, descErr
	}

	if descErr != nil {
		return nil, descErr
	}

	resultWrapper := NewResourceDescriptionResult(nil, info, output, rtype, keys, outId)
	return resultWrapper, descErr
}
