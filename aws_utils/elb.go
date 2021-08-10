package aws_utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	elb "github.com/aws/aws-sdk-go/service/elbv2"
	log "github.com/sirupsen/logrus"
)

const (
	ELBDnsFilter = "dns-filter"
)

type LBDescriptionInput struct {
	DNSNames      map[string]bool
	dnsNamesFound map[string]bool
	searchType    string
	Region        *string
}

func NewLBDescriptionInputFromDNS(dnsNames []string, region string) *LBDescriptionInput {
	m := map[string]bool{}
	dnsNamesFound := map[string]bool{}
	for _, addr := range dnsNames {
		m[NewDNS(addr).GetNormalAddr()] = true
	}

	return &LBDescriptionInput{
		Region:        aws.String(region),
		DNSNames:      m,
		searchType:    ELBDnsFilter,
		dnsNamesFound: dnsNamesFound,
	}
}

func (i *LBDescriptionInput) GetAWSInput() *elb.DescribeLoadBalancersInput {
	return &elb.DescribeLoadBalancersInput{}
}

func (i *LBDescriptionInput) isMatching(l *elb.LoadBalancer) (bool, bool, error) {

	if i.searchType != ELBDnsFilter {
		return false, false, fmt.Errorf("search type not support for elb - %s onl - %s ", i.searchType, ELBDnsFilter)
	}

	keepSearching := true

	ldns := NewDNSAWS(l.DNSName)

	_, found := i.DNSNames[ldns.GetNormalAddr()]

	if found {
		i.dnsNamesFound[ldns.GetNormalAddr()] = true
	}

	// if all dns names already found then there is no need to continue searching
	if len(i.DNSNames) == len(i.dnsNamesFound) {
		keepSearching = false
	}

	return found, keepSearching, nil
}

type LBDescriptionOutput struct {
	LoadBalancers []*elb.LoadBalancer
}

func (d *AWSResourceDescriber) lbclient(region string) *elb.ELBV2 {
	if d.elbClient == nil {
		d.elbClient = elb.New(d.session, aws.NewConfig().WithRegion(region))
		return d.elbClient
	}
	return d.elbClient
}

func (d *AWSResourceDescriber) describeLB(i *LBDescriptionInput) (*LBDescriptionOutput, error) {

	resultOutput := &LBDescriptionOutput{
		LoadBalancers: []*elb.LoadBalancer{},
	}

	var matchErrs error

	c := d.lbclient(*i.Region)

	input := i.GetAWSInput()

	pageCount := 0

	err := c.DescribeLoadBalancersPages(input, func(page *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {

		pageCount++
		keepSearching := false
		if page.LoadBalancers == nil {
			return keepSearching
		}
		for _, l := range page.LoadBalancers {

			isMatch, keepSearching, err := i.isMatching(l)

			lg := log.WithFields(log.Fields{
				"lb":            *l.LoadBalancerName,
				"currentPage":   pageCount,
				"isMatch":       isMatch,
				"keepSearching": keepSearching,
				"isLastPage":    lastPage,
			})

			if err != nil {
				matchErrs = err
				lg.WithError(err).Debug("error while check if dns match")
			}

			if isMatch {
				lg.Debug("found match")
				resultOutput.LoadBalancers = append(resultOutput.LoadBalancers, l)
			}

			if !keepSearching || err != nil {
				lg.Debug("stopping search for lb dns name")
				return keepSearching
			}
		}

		return keepSearching
	})

	if err != nil {
		log.WithError(err).Error("failed getting elb description")
		return resultOutput, err
	}

	if matchErrs != nil {
		return resultOutput, matchErrs
	}

	return resultOutput, err
}
