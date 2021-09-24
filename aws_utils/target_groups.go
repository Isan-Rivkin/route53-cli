package aws_utils

import (
	"sort"
	"strings"

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
	TGToTargets map[string][]*TGTargetOutput
}

// GetOutputID is describing the resources in the query in a unique way
// e.g if there are 3 instances their id will be unique and always consistent
// used to identify cache
func (tgo *TGDescriptionOutput) GetOutputID() string {
	arns := tgo.GetKeys()[ARNAttr]
	sortedArns := sort.StringSlice(arns)
	return strings.Join(sortedArns, ",")
}

func (tgo *TGDescriptionOutput) GetKeys() map[ResourceKey][]string {
	result := map[ResourceKey][]string{}
	var arns []string
	for _, tg := range tgo.TargetGroups {
		arns = append(arns, aws.StringValue(tg.TargetGroupArn))
	}

	result[ARNAttr] = arns

	if len(arns) > 0 {
		result[RegionAttr] = []string{NewDefaultResourceIdentifier().InferRegionFromResourceARN(arns[0])}
	}

	return result
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

		output.TGToTargets = map[string][]*TGTargetOutput{}

		for _, tg := range awsOutput.TargetGroups {
			targetInput := NewTargetHealthDescTGInputFromTargetGroupArn(*tg.TargetGroupArn, i.Region)
			targetsHealthOutput, err := d.describeTargetsHealth(targetInput)

			if err != nil {
				log.WithError(err).Error("failed describing target group targets health")
				continue
			}

			output.TGToTargets[*tg.TargetGroupArn] = targetsHealthOutput.Targets
		}
	}

	return output, nil
}
