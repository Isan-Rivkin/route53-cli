package aws_utils

import (
	"errors"
	"fmt"

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
	URLGenerators   = map[string]func(dnsIdentifier string, r *route53.ResourceRecordSet) string{
		ALBDns:        GetLBWebURL,
		NLBDns:        GetLBWebURL,
		CloudFrontDns: GetCloudFrontWebURL,
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
			return g(dnsType, r), nil
		}
		e = errors.New("ErrNotSupportedGenerator")
	}
	return "", e
}

func getRegionFromLBDNSName(dnsIdentifier string, r *route53.ResourceRecordSet) string {
	record := *r.AliasTarget.DNSName
	record = strings.TrimRight(record, ".")

	region := *r.Region
	// parse region
	splitted := strings.Split(record, ".")

	if dnsIdentifier == ALBDns || dnsIdentifier == ELBDns {
		// .us-east-1.elb.amazonaws.com
		if len(splitted) >= 4 {
			region = splitted[len(splitted)-4]
		}

	} else if dnsIdentifier == NLBDns {

		// .elb.us-east-1.amazonaws.com
		if len(splitted) >= 3 {
			region = splitted[len(splitted)-3]
		}
	}
	return region
}

// https://console.aws.amazon.com/ec2/v2/home?region=us-east-1#LoadBalancers:search=some-alb-name;sort=loadBalancerName
func GetLBWebURL(dnsIdentifier string, r *route53.ResourceRecordSet) string {
	record := *r.AliasTarget.DNSName
	record = strings.TrimRight(record, ".")

	searchQuery := record
	// parse region
	region := getRegionFromLBDNSName(dnsIdentifier, r)

	if strings.HasPrefix(searchQuery, "dualstack.") {
		searchQuery = strings.TrimLeft(searchQuery, "dualstack.")
	}
	return fmt.Sprintf("https://console.aws.amazon.com/ec2/v2/home?region=%s#LoadBalancers:search=%s;sort=loadBalancerName", region, searchQuery)
}

// https://console.aws.amazon.com/cloudfront/home?region=us-west-2#
// todo:: extract the record id and point to specific resource, cant search via url this is the general ui
// the region doesn't matter it's global
func GetCloudFrontWebURL(dnsIdentifier string, r *route53.ResourceRecordSet) string {
	return fmt.Sprintf("https://console.aws.amazon.com/cloudfront/home?region=%s#", "us-east-1")
}
