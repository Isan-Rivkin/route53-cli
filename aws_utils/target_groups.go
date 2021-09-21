package aws_utils

import (
	"github.com/aws/aws-sdk-go/aws"
	elb "github.com/aws/aws-sdk-go/service/elbv2"
	log "github.com/sirupsen/logrus"
)

type TGDescriptionInput struct {
	LoadBalancerArn string
	Region          string
	WithTargets     bool
}

func NewTGDescriptionInputFromELBArn(arn, region string, withTargets bool) *TGDescriptionInput {
	return &TGDescriptionInput{
		LoadBalancerArn: arn,
		Region:          region,
		WithTargets:     withTargets,
	}
}
func (i *TGDescriptionInput) GetAWSInput() *elb.DescribeTargetGroupsInput {
	return &elb.DescribeTargetGroupsInput{
		LoadBalancerArn: &i.LoadBalancerArn,
	}
}

type TGDescriptionOutput struct {
	TargetGroups []*elb.TargetGroup
	// key: target group arn, value: targets description for the tg (i.e ec2 instances)
	// optional, usually will be set if WithTargets:true in TGDescriptionInput
	TGToTargets map[string][]*elb.TargetHealthDescription
}

// this is the way to infer the instances and other targets in a target group from a target group arn.
// the describe of a target group does not provide the info about the targets, only that.
// https://stackoverflow.com/questions/64235538/how-to-get-target-group-of-running-ec2-instances-using-aws-apis
func (d *AWSResourceDescriber) describeTargetHealth(tgArn, region string) ([]*elb.TargetHealthDescription, error) {
	c := d.lbclient(region)

	input := &elb.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(tgArn),
	}

	output, err := c.DescribeTargetHealth(input)

	return output.TargetHealthDescriptions, err
}

func (d *AWSResourceDescriber) describeTG(i *TGDescriptionInput) (*TGDescriptionOutput, error) {
	c := d.lbclient(i.Region)
	awsInput := i.GetAWSInput()

	awsOutput, err := c.DescribeTargetGroups(awsInput)

	if err != nil {
		return nil, err
	}

	output := &TGDescriptionOutput{
		TargetGroups: awsOutput.TargetGroups,
	}

	// attach targets information such as ec2 id's
	if i.WithTargets {

		output.TGToTargets = map[string][]*elb.TargetHealthDescription{}

		for _, tg := range awsOutput.TargetGroups {

			targetsHealth, err := d.describeTargetHealth(*tg.TargetGroupArn, i.Region)
			if err != nil {
				log.WithError(err).Error("failed describing target group targets health")
				continue
			}

			output.TGToTargets[*tg.TargetGroupArn] = targetsHealth
		}
	}

	return output, nil
}
