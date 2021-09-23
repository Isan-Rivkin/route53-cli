package expander

import (
	awsu "r53/aws_utils"
)

type Command string

const (
	R53ExpandCmd Command = "r53_expand"
)

type CommandResult struct {
	CmdType      Command
	ResourceType awsu.ResourceType
	Payload      interface{}
}

type CommandExecutor interface {
	Execute(payload interface{})
}
