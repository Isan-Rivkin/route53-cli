package aws_utils

import (
	"errors"
	"net"
	"strconv"
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

type TGTargetOutput struct {
	TargetHealth *elb.TargetHealthDescription
	// target could be lambda, ec2 or ip. use the helper methods to get the type of the target
	TargetType ResourceType
}

func NewTGTargetOutput(t *elb.TargetHealthDescription) *TGTargetOutput {
	target := &TGTargetOutput{
		TargetHealth: t,
	}

	// infer target type
	tid := *t.Target.Id
	// check if ec2
	if strings.HasPrefix(tid, "i-") {
		target.TargetType = EC2Type
	} else if strings.HasPrefix(tid, "arn:aws:lambda:") { // check if lambda
		target.TargetType = LambdaType
	} else { // check if ip
		addr := net.ParseIP(tid)
		if addr != nil {
			target.TargetType = RawIpType
		}
	}
	return target
}

func (o *TGTargetOutput) IsEC2() bool {
	return o.TargetType == EC2Type
}

func (o *TGTargetOutput) IsLambda() bool {
	return o.TargetType == LambdaType
}

func (o *TGTargetOutput) IsRawIp() bool {
	return o.TargetType == RawIpType
}

func (o *TGTargetOutput) GetHealthcheckPort() string {
	if o.TargetHealth.HealthCheckPort != nil {
		return *o.TargetHealth.HealthCheckPort
	}
	return ""
}

func (o *TGTargetOutput) IsTargetHealthy() bool {
	return o.TargetHealth.TargetHealth != nil &&
		o.TargetHealth.TargetHealth.State != nil &&
		*o.TargetHealth.TargetHealth.State == "healthy"
}

func (o *TGTargetOutput) GetTargetID() string {
	return *o.TargetHealth.Target.Id
}

// lambda has no target port and if the protocol is GENEVE then port is 6081
func (o *TGTargetOutput) GetTargetPort() (string, error) {
	if o.IsLambda() {
		return "", errors.New("NoPortForLambdaTarget")
	}
	strPort := strconv.FormatInt(*o.TargetHealth.Target.Port, 10)
	return strPort, nil
}

type TGDescriptionOutput struct {
	TargetGroups []*elb.TargetGroup
	// key: target group arn, value: targets description for the tg (i.e ec2 instances)
	// optional, usually will be set if WithTargets:true in TGDescriptionInput
	TGToTargets map[string][]*TGTargetOutput
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

		output.TGToTargets = map[string][]*TGTargetOutput{}

		for _, tg := range awsOutput.TargetGroups {

			targetsHealth, err := d.describeTargetHealth(*tg.TargetGroupArn, i.Region)
			if err != nil {
				log.WithError(err).Error("failed describing target group targets health")
				continue
			}

			var targets []*TGTargetOutput
			for _, th := range targetsHealth {
				targets = append(targets, NewTGTargetOutput(th))
			}
			output.TGToTargets[*tg.TargetGroupArn] = targets
		}
	}

	return output, nil
}
