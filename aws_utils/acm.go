package aws_utils

import (
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
)

// TODO:: add option to filter all certificates based on other parameters such as elb arn, domains etc
type ACMDescpInput struct {
	Region string
	// if set then will get specific certificate
	CertARN *string
}

func NewACMDescpInputFromARN(certArn string) *ACMDescpInput {
	region := NewDefaultResourceIdentifier().InferRegionFromResourceARN(certArn)
	return &ACMDescpInput{
		CertARN: &certArn,
		Region:  region,
	}
}

func (i *ACMDescpInput) ToAWSInput() *acm.DescribeCertificateInput {
	return &acm.DescribeCertificateInput{
		CertificateArn: i.CertARN,
	}
}

type ACMDescOutput struct {
	Output *acm.DescribeCertificateOutput
}

// list of all domains in the certificate including the fqdn main domain
func (o *ACMDescOutput) GetAllDomains() []string {
	return aws.StringValueSlice(o.Output.Certificate.SubjectAlternativeNames)
}

func (o *ACMDescOutput) GetCertStatus() string {
	return aws.StringValue(o.Output.Certificate.Status)
}

func (o *ACMDescOutput) IsCertExpired() bool {
	now := time.Now()
	expiry := o.GetCertExpiryDate()
	return now.After(expiry)
}

func (o *ACMDescOutput) GetCertExpiryDate() time.Time {
	return aws.TimeValue(o.Output.Certificate.NotAfter)
}
func (o *ACMDescOutput) GetCertUsers() []string {
	return aws.StringValueSlice(o.Output.Certificate.InUseBy)
}

func (o *ACMDescOutput) IsCertInUse() bool {
	return len(o.GetCertUsers()) > 0
}

func (o *ACMDescOutput) GetAllDomainsStatus() []*acm.DomainValidation {
	return o.Output.Certificate.DomainValidationOptions
}

// checks if all the domains in the ceritificate are valid
func (o *ACMDescOutput) IsAllDomainsValidated() bool {
	for _, dv := range o.GetAllDomainsStatus() {
		if aws.StringValue(dv.ValidationStatus) != "SUCCESS" {
			return false
		}
	}
	return true
}

// GetOutputID is describing the resources in the query in a unique way
// e.g if there are 3 instances their id will be unique and always consistent
// used to identify cache
func (o *ACMDescOutput) GetOutputID() string {
	arns := o.GetKeys()[ARNAttr]
	sortedArns := sort.StringSlice(arns)
	return strings.Join(sortedArns, ",")
}

func (o *ACMDescOutput) GetKeys() map[ResourceKey][]string {

	result := map[ResourceKey][]string{}

	certArn := aws.StringValue(o.Output.Certificate.CertificateArn)
	result[ARNAttr] = []string{certArn}
	result[RegionAttr] = []string{NewDefaultResourceIdentifier().InferRegionFromResourceARN(certArn)}

	return result
}

func (d *AWSResourceDescriber) acmclient(region string) *acm.ACM {
	if d.acmClient == nil {
		d.acmClient = acm.New(d.session, aws.NewConfig().WithRegion(region))
		return d.acmClient
	}
	return d.acmClient
}

func (d *AWSResourceDescriber) describeCertificate(i *ACMDescpInput) (*ACMDescOutput, error) {

	input := i.ToAWSInput()
	c := d.acmclient(i.Region)

	out, err := c.DescribeCertificate(input)

	if err != nil {
		return nil, err
	}

	return &ACMDescOutput{Output: out}, nil
}
