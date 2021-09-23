package expander

import (
	"fmt"
	awsu "r53/aws_utils"
	"r53/cliui/abstracts"

	route53 "github.com/aws/aws-sdk-go/service/route53"
)

type R53ExpanderCmd struct {
	controller        Controller
	resourceIdenfiner awsu.AWSResourceIdentifier
	resourceDescriber awsu.AWSResourceDescriber
}

func NewR53ExpanderCmd(c Controller, identifier awsu.AWSResourceIdentifier, describer awsu.AWSResourceDescriber) CommandExecutor {
	return &R53ExpanderCmd{
		controller:        c,
		resourceIdenfiner: identifier,
		resourceDescriber: describer,
	}
}

func (c *R53ExpanderCmd) Execute(payload interface{}) {
	selection := payload.(*abstracts.TableSelectionResult)

	aliasRecord := selection.SelectedReference

	types, err := c.resourceIdenfiner.InferTypeFromDNS(aliasRecord)

	if err != nil {
		c.controller.SubmitResult(&CommandResult{
			CmdType: R53ExpandCmd,
		}, err)
		return
	}

	theType := types[0]

	region := c.resourceIdenfiner.InferRegionFromDNS(
		&route53.ResourceRecordSet{AliasTarget: &route53.AliasTarget{
			DNSName: &aliasRecord,
		}})

	// evaluate the r53 record -> resource type to expand
	var input interface{}

	switch theType {
	case awsu.ALBOrCLBType:
	case awsu.ELBType:
		input = awsu.NewLBDescriptionInputFromDNS([]string{aliasRecord}, region)
	default:
		c.controller.SubmitResult(&CommandResult{
			CmdType: R53ExpandCmd,
		}, fmt.Errorf("resource type not supported %s", theType))

		return
	}

	// describe the resource

	result, err := c.resourceDescriber.Describe(theType, input)

	if err != nil {
		c.controller.SubmitResult(&CommandResult{
			CmdType: R53ExpandCmd,
		}, fmt.Errorf("failed describing resource %s", theType))
		return
	}

	// success submit
	c.controller.SubmitResult(&CommandResult{
		ResourceType: theType,
		Payload:      result,
		CmdType:      R53ExpandCmd,
	}, nil)
}
