package aws_utils

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	route53 "github.com/aws/aws-sdk-go/service/route53"
)

type ZoneID string

func foo() int {
	return 1
}
func GetZones(sess *session.Session) ([]*route53.HostedZone, error) {
	r53 := route53.New(sess)
	r, err := r53.ListHostedZones(nil)
	if err != nil {
		return nil, err
	}
	return r.HostedZones, nil
}

func CreateListRecordSetInput(output *route53.ListResourceRecordSetsOutput, zoneId string) *route53.ListResourceRecordSetsInput {
	if output != nil {
		return &route53.ListResourceRecordSetsInput{
			HostedZoneId:          aws.String(zoneId), // Required
			StartRecordIdentifier: output.NextRecordIdentifier,
			StartRecordName:       output.NextRecordName,
			StartRecordType:       output.NextRecordType,
		}
	}
	return &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneId),
	}
}

// GetRecordSetsNoPage gets all records of a given zone ignoring pagination - i.e getting all
func GetRecordSetsNoPage(sess *session.Session, zoneID string) ([]*route53.ResourceRecordSet, error) {
	result := []*route53.ResourceRecordSet{}
	r53 := route53.New(sess)
	var output *route53.ListResourceRecordSetsOutput = nil
	for {
		listParams := CreateListRecordSetInput(output, zoneID)
		response, err := r53.ListResourceRecordSets(listParams)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		output = response
		result = append(result, response.ResourceRecordSets...)
		if *response.IsTruncated == false {
			break
		}
	}
	return result, nil
}

/* Route53 recods */
type R53RecordSetsResource struct {
	Records []*route53.ResourceRecordSet
}

func (i R53RecordSetsResource) Id() string {
	return "r53-records-set"
}

type R53RecordSetsHandler struct {
	RecordsResource *R53RecordSetsResource
}

func (h R53RecordSetsHandler) Id() string {
	return "r53-records-set"
}
func (h *R53RecordSetsHandler) Get() Resource {
	return h.RecordsResource
}
func (h *R53RecordSetsHandler) Fetch(config Conf, sess *session.Session, callback ResCb) {
	zoneId := config["zone-id"]
	if zoneId == nil {
		err := errors.New("missing zone-id")
		callback(err, nil)
	} else {
		response, err := GetRecordSetsNoPage(sess, zoneId.(string))
		if err != nil {
			callback(err, nil)
		} else {
			resource := &R53RecordSetsResource{Records: response}
			h.RecordsResource = resource
			callback(nil, resource)
		}
	}
}

/* Route53 zones */

type R53ZonesResource struct {
	Zones []*route53.HostedZone
}

func (i R53ZonesResource) Id() string {
	return "r53-zones"
}

type R53ZonesHandler struct {
	ZonesResource *R53ZonesResource
}

func (h R53ZonesHandler) Id() string {
	return "r53-zones"
}
func (h *R53ZonesHandler) Get() Resource {
	return h.ZonesResource
}
func (h *R53ZonesHandler) Fetch(config Conf, sess *session.Session, callback ResCb) {
	zones, err := GetZones(sess)
	if err != nil {
		callback(err, nil)
	} else {
		resource := &R53ZonesResource{Zones: zones}
		h.ZonesResource = resource
		callback(nil, h.ZonesResource)
	}
}

/** Route53 Zones and RecordSet*/

type R53ZonesWithRecordsResource struct {
	Zones      []*route53.HostedZone
	RecordsMap map[ZoneID][]*route53.ResourceRecordSet
}

func (i R53ZonesWithRecordsResource) Id() string {
	return "r53-zones-and-records-set"
}

type R53ZonesWithRecordsHandler struct {
	R53ZonesWithRecordsResource *R53ZonesWithRecordsResource
}

func (i R53ZonesWithRecordsHandler) Id() string {
	return "r53-zones-and-records-set"
}
func (h *R53ZonesWithRecordsHandler) Get() Resource {
	return h.R53ZonesWithRecordsResource
}
func (h *R53ZonesWithRecordsHandler) Fetch(config Conf, sess *session.Session, callback ResCb) {

}
