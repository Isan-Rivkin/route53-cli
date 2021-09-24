package aws_utils

import (
	"fmt"
)

type InputDescriber struct {
	Inputs []interface{}
}

func NewInput(input interface{}) *InputDescriber {
	var inputs []interface{}

	if IsSlice(input) {
		for _, i := range input.([]interface{}) {
			inputs = append(inputs, i)
		}
	} else {
		inputs = append(inputs, input)
	}

	return &InputDescriber{
		Inputs: inputs,
	}
}

func (id *InputDescriber) IsMultipleInputs() bool {
	return len(id.Inputs) > 1
}

// map of functions to create input for description from concrete resource
var describeInputBuilder = map[ResourceType]func(from Resource, to ResourceType) (*InputDescriber, error){
	ALBOrCLBType:        elbInputBuilder,
	ELBType:             elbInputBuilder,
	TargetGroupType:     targetGroupTypeInputBuilder,
	HealthTargetsTGType: healthTargetsTgTypeInputBuilder,
	ELBListenersType:    elbListenersTypeInputBuilder,
	ACMType:             acmTypeInputBuilder,
	EC2Type:             ec2TypeInputBuilder,
}

// Create input based on buildInput method
func BuildResourceInputFromInput(from Resource, to ResourceType) (*InputDescriber, error) {
	return buildInput(from, to)
}

// build input from source resource to a target resource
// - validate there is potential path from -> to
// - convert from to a concrete resource
// - call the concrete input builder func to create the input
func buildInput(from Resource, to ResourceType) (*InputDescriber, error) {

	if !IsReachableFrom(from.Type(), to) {
		return nil, fmt.Errorf("NoPathForResources from: %s to: %s", from.Type(), to)
	}

	builder, exist := describeInputBuilder[to]

	if !exist {
		return nil, fmt.Errorf("NoResourceBuilder: %s", to)
	}

	input, err := builder(from, to)

	return input, err

}

// create input for description from concrete resource
func elbInputBuilder(from Resource, to ResourceType) (*InputDescriber, error) {

	dnsRes, exist := from.TryGetValue(DNSAttr)

	if !exist {
		return nil, fmt.Errorf("MustContainDNSAttributeAddress from %s", from.Type())
	}

	dnsAlias := dnsRes[0]

	regionRes, exist := from.TryGetValue(RegionAttr)

	if !exist {
		return nil, fmt.Errorf("MustContainRegionAttribute from %s", from.Type())
	}

	region := regionRes[0]

	input := NewLBDescriptionInputFromDNS([]string{dnsAlias}, region)

	return NewInput(input), nil
}

// create input for description from concrete resource
func targetGroupTypeInputBuilder(from Resource, to ResourceType) (*InputDescriber, error) {

	arnsRes, exist := from.TryGetValue(ARNAttr)

	if !exist {
		return nil, fmt.Errorf("MustHaveArnsForLB from %s", from.Type())
	}

	arn := arnsRes[0]

	regionRes, exist := from.TryGetValue(RegionAttr)

	if !exist {
		return nil, fmt.Errorf("MustContainRegionAttribute from %s", from.Type())
	}

	region := regionRes[0]

	input := NewTGDescriptionInputFromELBArn(arn, region, false)

	return NewInput(input), nil
}

// create input for description from concrete resource
func healthTargetsTgTypeInputBuilder(from Resource, to ResourceType) (*InputDescriber, error) {
	var inputs []*TargetHealthDescTGInput

	arns, exist := from.TryGetValue(ARNAttr)

	if !exist {
		return nil, fmt.Errorf("MustContainArnAttribute from %s", from.Type())
	}

	regionRes, exist := from.TryGetValue(RegionAttr)

	if !exist {
		return nil, fmt.Errorf("MustContainRegionAttribute from %s", from.Type())
	}

	region := regionRes[0]

	for _, arn := range arns {
		input := NewTargetHealthDescTGInputFromTargetGroupArn(arn, region)
		inputs = append(inputs, input)
	}

	return NewInput(inputs), nil
}

// create input for description from concrete resource
func elbListenersTypeInputBuilder(from Resource, to ResourceType) (*InputDescriber, error) {
	arnsRes, exist := from.TryGetValue(ARNAttr)

	if !exist {
		return nil, fmt.Errorf("MustHaveArnsForLB from %s", from.Type())
	}

	arn := arnsRes[0]

	regionRes, exist := from.TryGetValue(RegionAttr)

	if !exist {
		return nil, fmt.Errorf("MustContainRegionAttribute from %s", from.Type())
	}

	region := regionRes[0]

	input := NewLBListenersDescribeInputFromELB(arn, region)
	return NewInput(input), nil
}

// create input for description from concrete resource
func acmTypeInputBuilder(from Resource, to ResourceType) (*InputDescriber, error) {
	var inputs []*ACMDescpInput
	arnsRes, exist := from.TryGetValue(CertArnAttr)

	if !exist {
		return nil, fmt.Errorf("MustHaveCertArnAttr from %s", from.Type())
	}

	for _, c := range arnsRes {
		input := NewACMDescpInputFromARN(c)
		inputs = append(inputs, input)
	}

	return NewInput(inputs), nil
}

// create input for description from concrete resource
func ec2TypeInputBuilder(from Resource, to ResourceType) (*InputDescriber, error) {
	ec2Res, exist := from.TryGetValue(Ec2IDAttr)

	if !exist {
		return nil, fmt.Errorf("MustHaveEC2ID from %s", from.Type())
	}

	regionRes, exist := from.TryGetValue(RegionAttr)

	if !exist {
		return nil, fmt.Errorf("MustContainRegionAttribute from %s", from.Type())
	}

	region := regionRes[0]

	input := NewEC2InstanceDescInputFromInstanceIds(region, ec2Res)

	return NewInput(input), nil
}
