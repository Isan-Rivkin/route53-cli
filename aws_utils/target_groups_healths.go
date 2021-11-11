package aws_utils

import (
	"errors"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	elb "github.com/aws/aws-sdk-go/service/elbv2"
)

type TargetHealthDescTGInput struct {
	TargetGroupArn string
	Region         string
}

func NewTargetHealthDescTGInputFromTargetGroupArn(arn, region string) *TargetHealthDescTGInput {
	return &TargetHealthDescTGInput{
		TargetGroupArn: arn,
		Region:         region,
	}
}

type TargetsHealthDescTGOutput struct {
	Region  string
	Targets []*TGTargetOutput
}

// GetOutputID is describing the resources in the query in a unique way
// e.g if there are 3 instances their id will be unique and always consistent
// used to identify cache
func (targets *TargetsHealthDescTGOutput) GetOutputID() string {
	arns := targets.GetKeys()[ARNAttr]
	sortedArns := sort.StringSlice(arns)
	return strings.Join(sortedArns, ",")
}

func (targets *TargetsHealthDescTGOutput) GetKeys() map[ResourceKey][]string {

	result := map[ResourceKey][]string{}
	var ids []string
	var ec2Ids []string

	for _, t := range targets.Targets {
		ids = append(ids, t.GetTargetID())

		if t.IsEC2() {
			ec2Ids = append(ec2Ids, t.GetTargetID())
		}
	}

	result[ARNAttr] = ids
	result[Ec2IDAttr] = ec2Ids

	if len(ids) > 0 {
		result[RegionAttr] = []string{targets.Region}
	}

	return result
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

// this is the way to infer the instances and other targets in a target group from a target group arn.
// the describe of a target group does not provide the info about the targets, only that.
// https://stackoverflow.com/questions/64235538/how-to-get-target-group-of-running-ec2-instances-using-aws-apis
func (d *AWSResourceDescriber) describeTargetsHealth(i *TargetHealthDescTGInput) (*TargetsHealthDescTGOutput, error) {
	c := d.lbclient(i.Region)

	input := &elb.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(i.TargetGroupArn),
	}

	output, err := c.DescribeTargetHealth(input)

	var targets []*TGTargetOutput

	for _, th := range output.TargetHealthDescriptions {
		targets = append(targets, NewTGTargetOutput(th))
	}

	return &TargetsHealthDescTGOutput{Targets: targets, Region: i.Region}, err
}
