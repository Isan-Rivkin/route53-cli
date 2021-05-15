/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"os"
	awsu "r53/aws_utils"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cfgFile string
var recordInput string
var awsProfile string
var hostedZoneDepth *int
var recursiveSearch *bool
var debug bool
var webUrl bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "r53 -r '*.some.dns.record.com'\n  r53 -r https://my.r53.website.com",
	Short: "Query route53 to get your dns record values",
	Long: `Query Route53 to get all sorts of information about a dns record. 
r53 will use your default AWS credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		if debug == true {
			log.SetLevel(log.DebugLevel)
		}

		if recordInput == "" {
			log.Error("query must not be empty use -r flag or --help")
			return
		}
		api := awsu.NewRoute53Api()
		result, err := api.GetRecordSetAliases(recordInput)
		if err != nil {
			log.WithError(err).Error("failed")
			return
		}

		if result == nil {
			log.Error("no result found")
			return
		}

		result.PrintTable(&awsu.PrintOptions{WebURL: webUrl})
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.r53.yaml)")
	rootCmd.PersistentFlags().StringVarP(&recordInput, "record", "r", "", "-r www.foo.app.com")
	//rootCmd.PersistentFlags().StringVarP(&awsProfile, "profile", "p", "default", "~/.aws/credentials chosen account")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Get verbose output about the process")
	// add web urls
	rootCmd.PersistentFlags().BoolVar(&webUrl, "url", true, "print url to the aws console that will display the resource")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".r53" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".r53")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
