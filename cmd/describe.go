/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	awsu "r53/aws_utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		log.Info("describe called")
		//name := "dualstack.some.dns.addr.from.r53."
		//describeELB(name, "us-east-1")
		describeELBDynamic()
	},
}

func describeELBDynamic() {
	// r53 query
	defaultDepth := 3
	result, err := GetR53Query(defaultDepth)

	if err != nil {
		log.WithError(err).Error("failed, potentially not authorized with aws")
		return
	}

	if len(result) != 1 {
		log.Error(fmt.Errorf("there is %d, only 1 hosted zone in the result sets is currently supported", len(result)))
		return
	}

	// filter the ELB's to describe

	identifier := awsu.NewDefaultResourceIdentifier()

	for _, record := range result[0].Records {
		types, err := identifier.InferFromRecordSet(record)

		dnsAlias := ""

		if record.AliasTarget != nil && record.AliasTarget.DNSName != nil {
			dnsAlias = *record.AliasTarget.DNSName
		}

		// those resource aliases values are not aws resources targets in r53, skip them, nothing to describe
		if err != nil {
			log.Debug("not aws resources can't describe, skipping: ", dnsAlias)
			continue
		}

		resourceType := types[0]

		// if were here this means we can describe the resource
		log.Debug("aws resource identified ", resourceType, dnsAlias)

		// get region
		region := identifier.InferRegionFromDNS(record)
		// describe the resdource if elb

		lg := log.WithFields(log.Fields{
			"dns":    dnsAlias,
			"type":   resourceType,
			"region": region,
		})

		if resourceType == awsu.ALBOrCLBType {

			lg.Info("try describe elb")

			elbOutput := describeELB(dnsAlias, region)

			if elbOutput != nil {
				lg.Info("try describe target groups")
				describeTG(*elbOutput.LoadBalancers[0].LoadBalancerArn, region)
			}
		}
	}

}

func describeTG(arn, region string) {
	describer := awsu.NewAWSResourceDescriber(awsProfile)
	withTargets := true
	input := awsu.NewTGDescriptionInputFromELBArn(arn, region, withTargets)

	out, err := describer.Describe(awsu.TargetGroupType, input)

	output := out.(*awsu.TGDescriptionOutput)

	if err != nil {
		log.WithError(err).Error("failed describing target group from load balancer from arn")
		return
	}

	log.Info("found target group ", len(output.TargetGroups))

	for _, tg := range output.TargetGroups {
		log.Info(tg.String())
		targets := output.TGToTargets[*tg.TargetGroupArn]
		log.Info("---> tg contains targets ", len(targets))
		for _, t := range targets {
			currently im at this stage: 
			get address -> get r53 records -> extract elbs -> describe target groups -> passing withTargets = true -> returns target group targerts descriptions i.e ec2 instances and stuff 
			its another api call and for some reason i got rate limited, some result worked for alb other didnt.
			action: figure out the rate limit, describe targets, turn them into ec2 targes (identify from id) and keep expending 
			log.Info(t.String())
		}
	}
}

func describeELB(name, region string) *awsu.LBDescriptionOutput {
	describer := awsu.NewAWSResourceDescriber(awsProfile)

	input := awsu.NewLBDescriptionInputFromDNS([]string{name}, region)

	output, err := describer.Describe(awsu.ELBType, input)

	out := output.(*awsu.LBDescriptionOutput)

	if err != nil {
		log.WithError(err).Error("failed describing")
		return nil
	}

	log.Info("found load balancers ", len(out.LoadBalancers))
	for _, lb := range out.LoadBalancers {
		log.Info(lb.String())
	}
	return out
}

func init() {
	rootCmd.AddCommand(describeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// describeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// describeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
