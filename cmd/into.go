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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	awsu "r53/aws_utils"
	"r53/cliui/abstracts"
	ui "r53/cliui/v1"
	expander "r53/expander"
)

// intoCmd represents the into command
var intoCmd = &cobra.Command{
	Use:   "into",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//dummy()
		v1()
	},
}

func GetRecordSetFromDNS(dns string, result []*awsu.GetRecordAliasesResult) *route53.ResourceRecordSet {
	for _, res := range result {
		if res.Records != nil {
			for _, record := range res.Records {
				if record.AliasTarget != nil && dns == aws.StringValue(record.AliasTarget.DNSName) {
					return record
				}
			}
		}
	}
	return nil
}

func v1() {
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

	// create channel to get ui even updates
	eventsController := make(chan *ui.AppEvent)

	app := ui.NewR53App(eventsController)

	app.SetR53RecordsQueryResult(result[0])

	// expander logic
	cache := expander.NewExpanderCache()
	describer := awsu.NewAWSResourceDescriber(awsProfile)
	expander := expander.NewExpander(describer, cache)

	// controller mock
	go func() {
		for {
			select {
			case event, keepOpen := <-eventsController:
				if keepOpen {
					tableSelection := event.EventPayload.(*abstracts.TableSelectionResult)
					rowCells := tableSelection.RowCells
					r53RecordVal := rowCells[tableSelection.ColSelected].Reference
					recordVal := r53RecordVal.(string)
					// trigger tree and expansion
					record := GetRecordSetFromDNS(recordVal, result)
					identifier := awsu.NewDefaultResourceIdentifier()
					region := identifier.InferRegionFromDNS(record)
					to, err := identifier.InferTypeFromDNS(recordVal)

					if err != nil || len(to) != 1 {
						panic(err)
					}

					resource := awsu.NewR53RecordSetDescriber(recordVal, region)
					connectedResources, _ := expander.Expand(resource, to[0])
					reachableResources := awsu.GetReachableResources(awsu.R53RecordSetType)
					var displayResources []string
					seen := map[string]bool{}

					for _, r := range reachableResources {
						if _, found := seen[string(r)]; !found {
							displayResources = append(displayResources, string(r))
							seen[string(r)] = true
						}
					}
					app.AddResourceExpansionTree(string(awsu.R53RecordSetType), displayResources)
					app.UpdateFocusTxtView(awsu.ALBOrCLBType, connectedResources)

				} else {
					fmt.Println("closing channel ")
					return
				}
			}
		}
	}()

	err = app.Run()

	if err != nil {
		panic(err)
	}

}

// func dummy() {
// 	defaultDepth := 3
// 	result, err := GetR53Query(defaultDepth)

// 	if err != nil {
// 		log.WithError(err).Error("failed, potentially not authorized with aws")
// 		return
// 	}

// 	if len(result) != 1 {
// 		log.Error(fmt.Errorf("there is %d, only 1 hosted zone in the result sets is currently supported", len(result)))
// 		return
// 	}

// 	// create channel to get ui even updates
// 	eventsController := make(chan *ui.AppEvent)

// 	app := ui.NewR53App(eventsController)

// 	app.SetR53RecordsQueryResult(result[0])

// 	go func() {
// 		for {
// 			select {
// 			case event, keepOpen := <-eventsController:
// 				if keepOpen {
// 					fmt.Println("got event! ", event.Type)
// 					tableSelection := event.EventPayload.(*abstracts.TableSelectionResult)
// 					rowCells := tableSelection.RowCells
// 					alisCol := 3
// 					fmt.Println(fmt.Sprintf("selected! row %d col %d txt %s", tableSelection.RowSelected, tableSelection.ColSelected, rowCells[alisCol].Text))
// 					fmt.Println("!!!!! ", rowCells[alisCol].Reference)
// 					// trigger tree and expansion
// 					app.AddResourceExpansionTree()

// 				} else {
// 					fmt.Println("closing channel ")
// 					return
// 				}
// 			}
// 		}
// 	}()

// 	err = app.Run()

// 	if err != nil {
// 		panic(err)
// 	}
// }
func init() {
	rootCmd.AddCommand(intoCmd)
}
