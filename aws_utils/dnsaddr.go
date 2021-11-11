package aws_utils

import (
	"strings"
)

const (
	removeLeft = "dualstack."
)

type DNSAddr struct {
	rawAddr  string
	cleanUrl string
}

func NewDNSAWS(addr *string) *DNSAddr {
	add := ""
	if addr != nil {
		add = *addr
	}
	return NewDNS(add)
}

func NewDNS(addr string) *DNSAddr {

	cleanUrl := strings.TrimRight(addr, ".")
	cleanUrl = strings.TrimLeft(cleanUrl, removeLeft)

	return &DNSAddr{
		rawAddr:  addr,
		cleanUrl: cleanUrl,
	}
}

func (addr *DNSAddr) GetNormalAddr() string {
	return addr.cleanUrl
}

func (addr *DNSAddr) IsEqual(other string) bool {
	other = strings.TrimRight(other, ".")
	other = strings.TrimLeft(other, removeLeft)
	return other == addr.cleanUrl
}
