package aws_utils_test

import (
	awsu "r53/aws_utils"
	"strings"
	"testing"
)

func TestStripRecord(t *testing.T) {
	cases := map[string]string{
		"ronen": "ronen",
		"https://" + "\052" + ".example.com/query?param1=foo&param2=bar&param3=baz": "*.example.com",
		"*.example.com/query?param1=foo&param2=bar&param3=baz":                      "*.example.com",
		"https://example.com/aaaa":                                                  "example.com",
		"http://*.example.com/sss":                                                  "*.example.com",
		"https://a.foo.us-east-1.int.example.io":                                    "a.foo.us-east-1.int.example.io",
		"*.a.foo.us-east-1.int.example.io":                                          "*.a.foo.us-east-1.int.example.io",
	}
	for test, expected := range cases {
		strippedList, err := awsu.StripRecord(test)
		stripped := strings.Join(strippedList, ".")
		if err != nil || stripped != expected {
			t.Fatalf("error striping record %s parsed=[%s] expected=[%s]", err, stripped, expected)
		}
	}
}

func GetAllOptionsForZoneNameTest(t *testing.T) {
	cases := map[string]string{
		"a":                                    "a",
		"*.a":                                  "a",
		"a.b":                                  "a.b,b",
		"*.a.b":                                "a.b.c,b.c,c",
		"*.a.b.c.d":                            "a.b.c.d,b.c.d,c.d,d",
		"a.b.c.d.e":                            "a.b.c.d.e,b.c.d.e,c.d.e,d.e,e",
		"us-east-1.a-b.foo-goo.int.example.io": "us-east-1.a-b.foo-goo.int.example.io,a-b.foo-goo.int.example.io,a-b.foo-goo.int.example.io,foo-goo.int.example.io,int.example.io,example.io,io",
	}
	for test, expected := range cases {
		r, err := awsu.NewRecordName(test)
		if err != nil {
			t.Fatalf("failed to create record name from test %s", err)
		}
		strippedList, err := r.GetAllOptionsForZoneName()
		stripped := strings.Join(strippedList, ",")
		if err != nil || stripped != expected {
			t.Fatalf("error striping record %s parsed=[%s] expected=[%s]", err, stripped, expected)
		}
	}
}
