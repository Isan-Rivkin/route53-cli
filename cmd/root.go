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
var debug *bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "r53",
	Short: "Query for Route53",
	Long:  `Query for Route53`,
	Run: func(cmd *cobra.Command, args []string) {
		if recordInput == "" {
			return
		}
		api := awsu.NewRoute53Api()
		result, err := api.GetRecordSetAliases(recordInput)
		if err != nil {
			log.WithError(err).Error("failed")
		}
		if result == nil {
			log.Error("no result found")
		}
		result.PrintTable()
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
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.r53.yaml)")
	rootCmd.PersistentFlags().StringVarP(&recordInput, "r", "r", "", "-d *.foo.com")
	rootCmd.PersistentFlags().StringVarP(&awsProfile, "profile", "p", "default", "~/.aws/credentials chosen account")
	hostedZoneDepth = rootCmd.PersistentFlags().IntP("depth", "d", 0, "hosted zone depth to look at (default is recursive)")
	recursiveSearch = rootCmd.PersistentFlags().BoolP("recursive", "", false, "recursive search by default is false expecting exact match on record")
	debug = rootCmd.PersistentFlags().BoolP("debug", "", true, "debug output")
	if *debug == true {
		log.SetLevel(log.DebugLevel)
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
