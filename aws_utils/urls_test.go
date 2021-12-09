package aws_utils_test

import (
	awsu "r53/aws_utils"
	"testing"

	route53 "github.com/aws/aws-sdk-go/service/route53"
	"github.com/stretchr/testify/assert"
)

func newRecSet(dns, region string) *route53.ResourceRecordSet {
	return &route53.ResourceRecordSet{
		Region: &region,
		AliasTarget: &route53.AliasTarget{
			DNSName: &dns,
		},
	}
}

func TestGenerateWebURL(t *testing.T) {
	//r.AliasTarget.DNSName
	cases := []struct {
		Record         *route53.ResourceRecordSet
		ShouldError    bool
		ExpectedOutput string
	}{
		{
			Record:         newRecSet("", ""),
			ShouldError:    true,
			ExpectedOutput: "",
		},
		{
			Record:         newRecSet("dualstack.my-production-alb.eu-east-1.elb.amazonaws.com.", "eu-east-1"),
			ShouldError:    false,
			ExpectedOutput: "https://console.aws.amazon.com/ec2/v2/home?region=eu-east-1#LoadBalancers:search=my-production-alb.eu-east-1.elb.amazonaws.com;sort=loadBalancerName",
		},
		{
			Record:         newRecSet("some-rec-ok-eu-east-1.similarweb.com.", "eu-east-1"),
			ShouldError:    true,
			ExpectedOutput: "",
		},
	}

	for _, c := range cases {
		url, err := awsu.GenerateWebURL(c.Record)

		isError := err != nil

		assert.Equal(t, c.ShouldError, isError, "error comparison dont match")
		assert.Equal(t, url, c.ExpectedOutput, "url output not as expected")
	}
}
