/*
Copyright © 2020 Isan Rivkin isanrivkin@gmail.com

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
	"io/ioutil"
	"os"

	awsu "github.com/isan-rivkin/route53-cli/aws_utils"
	"github.com/isan-rivkin/route53-cli/sdk"

	v "github.com/isan-rivkin/cliversioner"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AppVersion = "0.4.3"
)

var cfgFile string
var recordInput string
var awsProfile string
var outputFile string
var recusiveSearchMaxDepth *int
var recursiveSearch *bool
var outputJson bool
var debug bool
var webUrl bool
var skipNSVerification bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "r53 -r '*.some.dns.record.com'\n  r53 -r https://my.r53.website.com",
	Version: AppVersion,
	Short:   "Query route53 to get your dns record values",
	Long: `Query Route53 to get all sorts of information about a dns record. 
r53 will use your default AWS credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteR53()
		VersionCheck()

	},
}

func VersionCheck() {
	var err error
	optoutVar := "R53_VERSION_CHECK"
	i := v.NewInput("route53-cli", "https://github.com/isan-rivkin", AppVersion, &optoutVar)
	out, err := v.CheckVersion(i)

	if err != nil || out == nil {
		log.WithError(err).Debug("failed checking latest version from github.com")
		return
	}

	if out.Outdated {
		m := fmt.Sprintf("%s is not latest, %s, upgrade to %s", out.CurrentVersion, out.Message, out.LatestVersion)
		log.Warn(m)
	}
}

func ExecuteR53() {
	depth := *recusiveSearchMaxDepth
	if !*recursiveSearch {
		// default recurse 3 depth max
		depth = 3
	}

	input, err := sdk.NewInput(recordInput, awsProfile, debug, false, skipNSVerification, true, depth)

	if err != nil {
		log.WithError(err).Error("failed creating input")
		return
	}

	results, err := sdk.SearchR53(input)

	if err != nil {
		log.WithError(err).Error("failed searching")
		return
	}
	if outputJson || outputFile != "" {
		res, err := sdk.ToJSONOutput(results, true)
		if err != nil {
			log.WithError(err).Error("failed converting output to json")
			return
		}

		if outputFile != "" {
			if err = ioutil.WriteFile(outputFile, res, 0644); err != nil {
				log.WithError(err).Error("failed writting result to output file")
				return
			}
		} else {
			fmt.Println(string(res))
		}

	} else {
		for _, r := range results {
			r.PrintTable(&awsu.PrintOptions{WebURL: webUrl})
		}
	}

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
	rootCmd.PersistentFlags().StringVarP(&awsProfile, "profile", "p", "default", "~/.aws/credentials chosen account")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Get verbose output about the process")
	// add web urls
	rootCmd.PersistentFlags().BoolVar(&webUrl, "url", true, "print url to the aws console that will display the resource")
	rootCmd.PersistentFlags().BoolVar(&skipNSVerification, "ns-skip", false, "if set then nameservers will not be verified against the hosted zone result")
	R := false
	rootCmd.PersistentFlags().BoolVarP(&R, "recurse", "R", false, "if used then the tool will run recursively until all records have resolved")
	recursiveSearch = &R

	maxDepth := 3
	rootCmd.PersistentFlags().IntVarP(&maxDepth, "max-depth", "d", maxDepth, "if -R is used then specifies when to stop recursive search depth")
	recusiveSearchMaxDepth = &maxDepth
	// output as json result
	rootCmd.PersistentFlags().BoolVar(&outputJson, "output-json", false, "output the result as a json object")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output-file", "f", "", "save output as json to file path")
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
