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

package sdk

import (
	"fmt"

	awsu "github.com/isan-rivkin/route53-cli/aws_utils"

	log "github.com/sirupsen/logrus"
)

type AWSAuthInput struct {
	AWSProfile string
}

type Input struct {
	SkipNSVerification bool
	AwsAuth            AWSAuthInput
	Record             string
	RecursiveSearch    bool
	RecursiveMaxDepth  int
	Debug              bool
	MuteLogs           bool
}

func NewInput(recordInput, awsProfile string, debug, muteLogs bool, skipNSVerification, recursiveSearch bool, recursiveMaxDepth int) (Input, error) {
	i := Input{}

	if recordInput == "" {
		return i, fmt.Errorf("record input to search is empty")
	}

	i.SkipNSVerification = skipNSVerification
	i.AwsAuth = AWSAuthInput{AWSProfile: awsProfile}
	i.Record = recordInput
	i.Debug = debug
	i.MuteLogs = muteLogs
	i.RecursiveSearch = recursiveSearch
	i.RecursiveMaxDepth = recursiveMaxDepth

	if i.Debug && !i.MuteLogs {
		log.SetLevel(log.DebugLevel)
	}
	if i.MuteLogs {
		log.SetLevel(log.PanicLevel)
	}

	if skipNSVerification {
		log.Warn("skipping nameserver verification, possibly inccorect result, not recomended.")
	}

	log.WithField("profile", awsProfile).Info("using aws environment session")

	return i, nil
}

func SearchR53(in Input) ([]*awsu.GetRecordAliasesResult, error) {
	api := awsu.NewRoute53Api(in.AwsAuth.AWSProfile)

	results, err := api.GetRecordSetAliasesRecursive(in.RecursiveMaxDepth, in.Record, in.SkipNSVerification, nil)

	if err != nil {
		return nil, err
	}
	return results, nil
}
