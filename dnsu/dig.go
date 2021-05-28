package dnsu

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lixiangzhong/dnsutil"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type NSRecord struct {
	Name string
	RR   dns.RR
}

type NSLookupResult struct {
	Records map[string]*NSRecord
}

func (nsr *NSLookupResult) IsRecordExist(record string) bool {
	r := record
	if !strings.HasSuffix(record, ".") {
		r += "."
	}

	_, exist := nsr.Records[r]
	return exist
}

func (nsr *NSLookupResult) AppendRecords(other *NSLookupResult) {

	if other != nil && len(other.Records) > 0 {

		if nsr.Records == nil {
			nsr.Records = map[string]*NSRecord{}
		}

		for name, rec := range other.Records {
			nsr.Records[name] = rec
		}
	}
}
func DigNSSearch(domain string) (*NSLookupResult, error) {

	var dig dnsutil.Dig

	msg, err := dig.GetMsg(dns.TypeNS, domain)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return newNSLookupResultFromDigMsg(msg)
}

func newNSLookupResultFromDigMsg(m *dns.Msg) (*NSLookupResult, error) {
	result := &NSLookupResult{
		Records: map[string]*NSRecord{},
	}

	if ans, err := praseRRFromDigMsg(m.Answer); err == nil {
		result.AppendRecords(ans)
	}
	if extra, err := praseRRFromDigMsg(m.Extra); err == nil {
		result.AppendRecords(extra)
	}
	if ns, err := praseRRFromDigMsg(m.Ns); err == nil {
		result.AppendRecords(ns)
	}

	if result.Records == nil || len(result.Records) == 0 {
		return nil, errors.New("ErrNoNSResultForHost")
	}
	return result, nil
}

func praseRRFromDigMsg(rr []dns.RR) (*NSLookupResult, error) {

	result := &NSLookupResult{
		Records: map[string]*NSRecord{},
	}

	if len(rr) == 0 {
		return result, errors.New("ErrNoRR")
	}

	for _, r := range rr {
		if t, ok := r.(*dns.NS); ok {

			log.WithFields(log.Fields{
				"type":   "NS",
				"result": t,
			}).Debug("found NS record")

			result.Records[t.Ns] = &NSRecord{
				Name: t.Ns,
				RR:   t,
			}
			// SOA for private hosted zone
		} else if t, ok := r.(*dns.SOA); ok {

			log.WithFields(log.Fields{
				"type":   "SOA",
				"result": t,
			}).Debug("found SOA record")

			result.Records[t.Ns] = &NSRecord{
				Name: t.Ns,
				RR:   t,
			}
		}
	}
	return result, nil
}

func DigTest3(domain string) {
	var dig dnsutil.Dig

	msg, err := dig.GetMsg(dns.TypeNS, domain)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(msg)
	// or
	fmt.Println("@@@@@@@@@@@@@@@@@")
	fmt.Println(msg.Question) //question section
	fmt.Println(msg.Answer)   //answer section.
	fmt.Println("ANSWER --- >> !!!!!!!!!!!!!!!!!!!")
	for _, r1 := range msg.Answer {
		fmt.Println(r1)
	}
	fmt.Println("ANSWER --- >>  !!!!!!!!!!!!!!!!!!!")
	fmt.Println(msg.Ns) //authority section.
	if t, ok := msg.Answer[0].(*dns.NS); ok {
		// do something with t.Txt
		fmt.Println("hi hi hi  ok :)", t.Ns)
	}
	fmt.Println("!!!!!!!!!!!!!!!!!!!")
	for _, r1 := range msg.Ns {
		fmt.Println(r1)
	}
	// if t, ok := msg.Ns[0].(*dns.SOA); ok {
	// 	// do something with t.Txt
	// 	fmt.Println("ok ok :)", t.Ns)
	// }

	fmt.Println("!!!!!!!!!!!!!!!!!!!")
	fmt.Println(msg.Extra) //additional section.

	fmt.Println("extra --- !!!!!!!!!!!!!!!!!!!")
	for _, r1 := range msg.Extra {
		fmt.Println(r1.Header().Name)
	}
	fmt.Println("extra --- !!!!!!!!!!!!!!!!!!!")
	fmt.Println("^^^^^^^^^^^^^^")
}
