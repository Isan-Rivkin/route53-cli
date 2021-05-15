package aws_utils

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"strings"

	route53 "github.com/aws/aws-sdk-go/service/route53"
)

const (
	// .us-east-1.elb.amazonaws.com
	ALBDns = "elb.amazonaws.com"
	CLBDns = "elb.amazonaws.com"
	// .elb.us-east-1.amazonaws.com
	NLBDns = "elb."
	// .elasticloadbalancing.us-east-2.amazonaws.com
	ELBDns = "elasticloadbalancing."
	// s3-website.af-south-1.amazonaws.com
	S3Suffx = "s3-website."
	// .us-east-1.vpce.amazonaws.com
	VPCEndpointDns = "vpce.amazonaws.com"
	CloudFrontDns  = "cloudfront.net"
	// .us-west-2.elasticbeanstalk.com
	ElasticBeansTalkDns = "elasticbeanstalk.com"
	AcceleratorApiDns   = "awsglobalaccelerator.com"
)

var (
	SupportedTarget = []string{ALBDns, CLBDns, NLBDns, ELBDns, S3Suffx, VPCEndpointDns, CloudFrontDns, ElasticBeansTalkDns, AcceleratorApiDns}
	URLGenerators   = map[string]func(r *route53.ResourceRecordSet) string{
		ALBDns: GetLBWebURL,
		NLBDns: GetLBWebURL,
	}
)

// checks if the value of a route53 record is an alias to a routable AWS resource, if true will return the resource type from the list
func CheckRoutableAWSTarget(r *route53.ResourceRecordSet) (string, bool) {
	if r.AliasTarget == nil || r.AliasTarget.DNSName == nil {
		return "", false
	}
	dns := *r.AliasTarget.DNSName

	for _, st := range SupportedTarget {
		if strings.Contains(dns, st) {
			return st, true
		}
	}
	return "", false
}

func GenerateWebURL(r *route53.ResourceRecordSet) (string, error) {
	e := errors.New("ErrNotRoutable")
	if dnsType, routable := CheckRoutableAWSTarget(r); routable {
		g, found := URLGenerators[dnsType]
		if found {
			return g(r), nil
		}
		e = errors.New("ErrNotSupportedGenerator")
	}
	return "", e
}

// https://console.aws.amazon.com/ec2/v2/home?region=us-east-1#LoadBalancers:search=web-prd-nomad-alb-ext-1356510359;sort=loadBalancerName
func GetLBWebURL(r *route53.ResourceRecordSet) string {
	record := *r.AliasTarget.DNSName
	record = strings.TrimRight(record, ".")
	log.Warn(r.Region)
	region := *r.Region
	searchQuery := record

	// parse region

	splitted := strings.Split(record, ".")
	// .us-east-1.elb.amazonaws.com
	if len(splitted) >= 4 {
		region = splitted[len(splitted)-4]
	}

	if strings.HasPrefix(searchQuery, "dualstack.") {
		searchQuery = strings.TrimLeft(searchQuery, "dualstack.")
	}

	return fmt.Sprintf("https://console.aws.amazon.com/ec2/v2/home?region=%s#LoadBalancers:search=%s;sort=loadBalancerName", region, searchQuery)
}
