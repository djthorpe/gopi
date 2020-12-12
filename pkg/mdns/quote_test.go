package mdns_test

import (
	"strconv"
	"testing"

	// Units
	mdns "github.com/djthorpe/gopi/v3/pkg/mdns"
)

var (
	quote_tests = []struct {
		unquoted, quoted string
		err              bool
	}{
		{"", "", false},
		{"abc", "abc", false},
		{"a-b", "a-b", false},
		{"a_b", "a_b", false},
		{"a b", "a\\ b", false},
		{"a\nb", "a\\nb", false},
		{"a\rb", "a\\rb", false},
		{"a\fb", "a\\fb", false},
		{"a\tb", "a\\tb", false},
		{"a@b", "a\\@b", false},
		{"a.b", "a\\.b", false},
		{"a`b", "a\\096b", false},
		{"Test’s Test", "Test\\226\\128\\153s\\ Test", false},
		{"Brother HL-3170CDW series", "Brother\\ HL-3170CDW\\ series", false},
		{"50-34-10-70.1 Backup", "50-34-10-70\\.1\\ Backup", false},
		{"fc:b6:d8:71:ac:4b@fe80::feb6:d8ff:fe71:ac4b", "fc:b6:d8:71:ac:4b\\@fe80::feb6:d8ff:fe71:ac4b", false},
		{"7c1c5acf-5a7e-d5bb-9593-e6d52e05f031", "7c1c5acf-5a7e-d5bb-9593-e6d52e05f031", false},
		{"Not\\", "Not\\\\", false},
		{"", "Not\\", true},
		{"Not\\045", "Not\\\\045", false},
	}
)

func Test_Quote_001(t *testing.T) {
	for i, test := range quote_tests {
		if test.err == true {
			continue
		}
		t.Logf("Test %v: %v => %v", i, strconv.Quote(test.unquoted), strconv.Quote(test.quoted))
		quoted := mdns.Quote(test.unquoted)
		if quoted != test.quoted {
			t.Errorf("Test %v: Expected %v but got %v", i, strconv.Quote(test.quoted), strconv.Quote(quoted))
		}
	}
}
func Test_Quote_002(t *testing.T) {
	for i, test := range quote_tests {
		t.Logf("Test %v: %v => %v", i, strconv.Quote(test.quoted), strconv.Quote(test.unquoted))
		if unquoted, err := mdns.Unquote(test.quoted); test.err {
			if err == nil {
				t.Errorf("Test %v: Expected error condition but got none", i)
			}
		} else if err != nil {
			t.Errorf("Test %v: Unexpected error condition: %v", i, err)
		} else if unquoted != test.unquoted {
			t.Errorf("Test %v: Expected %v but got %v", i, strconv.Quote(test.unquoted), strconv.Quote(unquoted))
		}
	}
}
