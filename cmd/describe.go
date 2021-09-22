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

	"github.com/aws/aws-sdk-go/aws"
	elb "github.com/aws/aws-sdk-go/service/elbv2"
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

				// target groups
				lg.Info("try describe target groups")
				tg := describeTG(*elbOutput.LoadBalancers[0].LoadBalancerArn, region)
				describeTargetGroupInstances(tg, region)

				// listeners
				// lg.Info("try describe elb listeners")
				// listeners := describeELBListeners(*elbOutput.LoadBalancers[0].LoadBalancerArn, region)
				// // certificates
				// lg.Info("try describe elb certificates from listeners")
				// for _, l := range listeners {
				// 	for _, c := range l.Certificates {
				// 		describeCertificates(*c.CertificateArn)
				// 	}
				// }
			}
		}
	}

}

func describeTargetGroupInstances(tg *awsu.TGDescriptionOutput, region string) {
	var ec2Ids []string

	for _, g := range tg.TargetGroups {
		arn := aws.StringValue(g.TargetGroupArn)
		targets := tg.TGToTargets[arn]
		for _, t := range targets {
			if t.IsEC2() {
				ec2Ids = append(ec2Ids, t.GetTargetID())
			}
		}
	}
	describeInstances(ec2Ids, region)
}

func describeInstances(ids []string, region string) {
	describer := awsu.NewAWSResourceDescriber(awsProfile)
	input := awsu.NewEC2InstanceDescInputFromInstanceIds(region, ids)

	out, err := describer.Describe(awsu.EC2Type, input)

	if err != nil {
		log.WithError(err).Error("failed describing ec2 instances from ids")
		return
	}

	output := out.(*awsu.EC2InstanceDescOutput)

	log.Info("tg has instances found ", len(output.GetAllInstances()))

	for _, inst := range output.GetAllInstances() {
		uptime, _ := inst.GetPrettyUptime()

		log.WithFields(log.Fields{
			"name":      inst.GetNameTag(),
			"id":        inst.ID(),
			"privateIp": inst.GetInstancePrivateIpV4(),
			"publicIp":  inst.GetInstancePublicIpV4(),
			"uptime":    uptime,
			"status":    inst.GetInstanceState(),
		}).Info("ec2 instance")
	}

}

func describeCertificates(certArn string) {
	describer := awsu.NewAWSResourceDescriber(awsProfile)
	input := awsu.NewACMDescpInputFromARN(certArn)

	out, err := describer.Describe(awsu.ACMType, input)

	if err != nil {
		log.WithError(err).Error("failed describing acm certificate from arn")
		return
	}

	output := out.(*awsu.ACMDescOutput)

	log.Info(output.GetAllDomains())
}

func describeELBListeners(arn, region string) []*elb.Listener {
	describer := awsu.NewAWSResourceDescriber(awsProfile)
	input := awsu.NewLBListenersDescribeInputFromELB(arn, region)
	out, err := describer.Describe(awsu.ELBListenersType, input)

	if err != nil {
		log.WithError(err).Error("failed describing listeners from load balancer from arn")
		return nil
	}

	output := out.(*awsu.LBListenersDescriptionOutput)

	listeners := output.GetListeners(arn)

	log.Info("listeners for elb found ", len(listeners))

	for _, l := range listeners {
		log.Info(l.String())
	}

	return listeners
}

func describeTG(arn, region string) *awsu.TGDescriptionOutput {
	describer := awsu.NewAWSResourceDescriber(awsProfile)
	withTargets := true
	input := awsu.NewTGDescriptionInputFromELBArn(arn, region, withTargets)

	out, err := describer.Describe(awsu.TargetGroupType, input)

	if err != nil {
		log.WithError(err).Error("failed describing target group from load balancer from arn")
		return nil
	}

	output := out.(*awsu.TGDescriptionOutput)

	log.Info("found target group ", len(output.TargetGroups))

	for _, tg := range output.TargetGroups {
		log.Info(tg.String())
		targets := output.TGToTargets[*tg.TargetGroupArn]
		log.Info("---> tg contains targets ", len(targets))
		for _, t := range targets {
			l := log.WithFields(log.Fields{
				"type":    t.TargetType,
				"id":      t.GetTargetID(),
				"healthy": t.IsTargetHealthy(),
			})
			l.Info("tg target")
		}
	}
	return output
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
