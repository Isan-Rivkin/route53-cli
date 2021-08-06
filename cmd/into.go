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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	ui "r53/cliui"
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

		s, err := ui.SelectR53RecordFromList(result[0])

		if err != nil {
			panic(err)
		}
		fmt.Println("selection = ", s)

		VersionCheck()
	},
}

func init() {
	rootCmd.AddCommand(intoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// intoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// intoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
