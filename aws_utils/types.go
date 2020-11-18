package aws_utils

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

// general configuration for handlers
type Conf map[string]interface{}

// result of a handler is ready
type ResCb func(error, Resource)

// Specific resource handler (for example: describe ec2 instances)
type Handler interface {
	Id() string
	Get() Resource
	Fetch(Conf, *session.Session, ResCb)
}

// Resource is a result of each handler
type Resource interface {
	Id() string
}

// Manages all Handlers for a specific group (i.e route53 is a group of many resources)
type GroupManager interface {
	Id() string
	GetHandler(string) Handler
	GetResource(string) Resource
}
