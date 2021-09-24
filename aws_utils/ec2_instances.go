package aws_utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EC2InstanceWrapper struct {
	Instance *ec2.Instance
}

func (e *EC2InstanceWrapper) ID() string {
	return aws.StringValue(e.Instance.InstanceId)
}

//pending, running, shutting-down, terminated, stopping, stopped
func (e *EC2InstanceWrapper) GetInstanceState() string {
	return aws.StringValue(e.Instance.State.Name)
}

func (e *EC2InstanceWrapper) GetInstancePublicIpV4() string {
	return aws.StringValue(e.Instance.PublicIpAddress)
}

func (e *EC2InstanceWrapper) GetInstancePrivateIpV4() string {
	return aws.StringValue(e.Instance.PrivateIpAddress)
}

func (e *EC2InstanceWrapper) GetInstanceIAMProfile() *ec2.IamInstanceProfile {
	return e.Instance.IamInstanceProfile
}

func (e *EC2InstanceWrapper) GetSecurityGroupIDs() []string {
	var sgids []string
	for _, sg := range e.Instance.SecurityGroups {
		sgids = append(sgids, aws.StringValue(sg.GroupId))
	}
	return sgids
}

func (e *EC2InstanceWrapper) GetTags() map[string]string {
	tags := map[string]string{}
	for _, t := range e.Instance.Tags {
		tags[aws.StringValue(t.Key)] = aws.StringValue(t.Value)
	}
	return tags
}

// returns instance uptime if running in a pretty format
func (e *EC2InstanceWrapper) GetPrettyUptime() (string, error) {

	if e.GetInstanceState() != "running" {
		return "", fmt.Errorf("InstanceNotRunningForUptime - %s ", e.GetInstanceState())
	}

	lunch := aws.TimeValue(e.Instance.LaunchTime)
	return GetPrettyUptime(lunch), nil
}

func (e *EC2InstanceWrapper) IsSpot() bool {
	return aws.StringValue(e.Instance.InstanceLifecycle) == "spot"
}

// describes the instance name instead of instance id
func (e *EC2InstanceWrapper) GetNameTag() string {
	tags := e.GetTags()
	return tags["Name"]
}

type EC2InstanceDescInput struct {
	Region string
	// optional, if not set default is all instances
	InstanceIds []*string
	// optional, maybe pre mature optimization, will filter based on attributes
	// https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#DescribeInstancesInput
	Filters []*ec2.Filter
}

func NewEC2InstanceDescInputFromInstanceIds(region string, ids []string) *EC2InstanceDescInput {
	return &EC2InstanceDescInput{
		Region:      region,
		InstanceIds: aws.StringSlice(ids),
	}
}

func (i *EC2InstanceDescInput) ToAWSInput() *ec2.DescribeInstancesInput {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: i.InstanceIds,
		Filters:     i.Filters,
	}
	return input
}

type EC2InstanceDescOutput struct {
	Region string
	// use helper methods here
	Output []*ec2.Reservation
}

func (o EC2InstanceDescOutput) GetAllInstances() []*EC2InstanceWrapper {
	var instances []*EC2InstanceWrapper
	for _, r := range o.Output {
		for _, inst := range r.Instances {
			instances = append(instances, &EC2InstanceWrapper{Instance: inst})
		}
	}
	return instances
}

// GetOutputID is describing the resources in the query in a unique way
// e.g if there are 3 instances their id will be unique and always consistent
// used to identify cache
func (o *EC2InstanceDescOutput) GetOutputID() string {
	arns := o.GetKeys()[ARNAttr]
	sortedArns := sort.StringSlice(arns)
	return strings.Join(sortedArns, ",")
}

func (o *EC2InstanceDescOutput) GetKeys() map[ResourceKey][]string {

	result := map[ResourceKey][]string{}
	var ids []string

	for _, i := range o.GetAllInstances() {
		ids = append(ids, i.ID())
	}

	result[ARNAttr] = ids

	if len(ids) > 0 {
		result[RegionAttr] = []string{o.Region}
	}

	return result
}

func (d *AWSResourceDescriber) ec2client(region string) *ec2.EC2 {
	if d.ec2Client == nil {
		d.ec2Client = ec2.New(d.session, aws.NewConfig().WithRegion(region))
		return d.ec2Client
	}
	return d.ec2Client
}

func (d *AWSResourceDescriber) describeEC2Instances(i *EC2InstanceDescInput) (*EC2InstanceDescOutput, error) {

	input := i.ToAWSInput()
	c := d.ec2client(i.Region)

	allOutput := []*ec2.Reservation{}

	err := c.DescribeInstancesPages(input, func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
		keepSearch := !lastPage

		allOutput = append(allOutput, page.Reservations...)

		return keepSearch
	})

	if err != nil {
		return nil, err
	}

	return &EC2InstanceDescOutput{Output: allOutput, Region: i.Region}, nil
}
