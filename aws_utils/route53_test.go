package aws_utils

import (
	"testing"
)

func TestHostedZones(t *testing.T) {
	sess := GetEnvSession()
	GetZones(sess)
}

func TestRecordSetsNoPage(t *testing.T) {
	sess := GetEnvSession()
	res, _ := GetRecordSetsNoPage(sess, "/hostedzone/Z1AEL1R23MWQEZ")
	if len(res) < 250 {
		t.Errorf("TestEc2InstancesResource(): %d < 250", len(res))
	}
}
