package cmds

import (
	"fmt"
	awsu "r53/aws_utils"
	"r53/cliui/abstracts"
	e "r53/expander"

	route53 "github.com/aws/aws-sdk-go/service/route53"
)

type R53ExpanderCmd struct {
	controller        e.Controller
	resourceIdenfiner awsu.AWSResourceIdentifier
	resourceDescriber awsu.AWSResourceDescriber
}

func NewR53ExpanderCmd(c e.Controller, identifier awsu.AWSResourceIdentifier, describer awsu.AWSResourceDescriber) e.CommandExecutor {
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
		c.controller.SubmitResult(&e.CommandResult{
			CmdType: e.R53ExpandCmd,
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
	case awsu.ELBType, awsu.ALBOrCLBType:
		input = awsu.NewLBDescriptionInputFromDNS([]string{aliasRecord}, region)
	default:
		c.controller.SubmitResult(&e.CommandResult{
			CmdType: e.R53ExpandCmd,
		}, fmt.Errorf("resource type not supported %s", theType))

		return
	}

	// describe the resource

	result, err := c.resourceDescriber.Describe(theType, input)

	if err != nil {
		c.controller.SubmitResult(&e.CommandResult{
			CmdType: e.R53ExpandCmd,
		}, fmt.Errorf("failed describing resource %s", theType))
		return
	}

	// success submit
	c.controller.SubmitResult(&e.CommandResult{
		ResourceType: theType,
		Payload:      result,
		CmdType:      e.R53ExpandCmd,
	}, nil)
}
