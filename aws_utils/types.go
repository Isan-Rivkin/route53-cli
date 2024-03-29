package aws_utils

import (
	"github.com/aws/aws-sdk-go/aws/session"
	route53 "github.com/aws/aws-sdk-go/service/route53"
)

type DNSRecordType string

const (
	WildCard string = "\\052"
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

type GetRecordAliasesResult struct {
	Records    []*route53.ResourceRecordSet
	HostedZone *route53.HostedZone
	Stream     RecordStream
	Region     string
}

type Route53Api interface {
	TestDNSAnswer(hostedZoneId, recordName, recordType string) (*route53.TestDNSAnswerOutput, error)
	GetHostedZonesFromDns(recordName string) ([]*route53.HostedZone, error)
	GetRecordSetAliases(recordName string, skipNSVerification bool) (*GetRecordAliasesResult, error)
	GetRecordSetAliasesRecursive(maxDepth int, recordName string, skipNSVerification bool, checkedRecord map[string]bool) ([]*GetRecordAliasesResult, error)
	GetRegion() string
	GetNameservers(recordName string) ([]string, error)
	GetHZNameservers(hzId string) ([]string, error)
	IsNSMatch(ns1, ns2 []string) bool
}

type Route53Manager struct {
	session     *session.Session
	r53client   *route53.Route53
	nameservers map[string][]string
}

type RecordStream interface {
	GetAllOptionsForZoneName() ([]string, error)
	IsEqual(other string) bool
	HasWildCard() bool
	GetWithWildCard() (string, error)
	GetParsedURL() string
}

type RecordName struct {
	hasWildCard bool
	rawURL      string
	parsedURL   string
	splittedURL []string
}
