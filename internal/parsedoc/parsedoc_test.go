package parsedoc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDoc(t *testing.T) {
	tests := map[string]struct {
		input      string
		outHeaders map[string]string
		outMessage map[string]string
		err        error
	}{
		"empty doc": {
			input:      "",
			outHeaders: map[string]string{},
			outMessage: map[string]string{},
			err:        nil,
		},
		"no headers": {
			input:      "\nk1:v1\nk2: v2",
			outHeaders: map[string]string{},
			outMessage: map[string]string{"k1": "v1", "k2": "v2"},
			err:        nil,
		},
		"all headers": {
			input: "anything_else: AE\ncustomer_code: C\nhost: H\nindex: I\nsource: S\nsource_environment: SE\nsourcetype: ST\ntype: T\n\nk1:v1",
			outHeaders: map[string]string{
				"customer_code":      "C",
				"host":               "H",
				"index":              "I",
				"source":             "S",
				"source_environment": "SE",
				"sourcetype":         "ST",
				"type":               "T",
				"anything_else":      "AE",
			},
			outMessage: map[string]string{"k1": "v1"},
			err:        nil,
		},
		"edge cases": {
			input:      "X-key:123\n\nk0\nk1:v1\nk2: v2\nk3:\n k4:no\nk5:5\nk6:",
			outHeaders: map[string]string{"X-key": "123"},
			outMessage: map[string]string{"k1": "v1", "k2": "v2", "k3": "", "k5": "5", "k6": ""},
			err:        nil,
		},
		"orion Informational w/out type": {
			input:      "\nk1:v1\nk2: v2\nvendor_severity: Informational",
			outHeaders: map[string]string{},
			outMessage: map[string]string{"k1": "v1", "k2": "v2", "vendor_severity": "Informational"},
			err:        nil,
		},
		"orion Informational": {
			input:      "type: STI::CE::Message::RawEventOrion\n\nk1:v1\nk2: v2\nvendor_severity: Informational",
			outHeaders: map[string]string{"type": "STI::CE::Message::RawEventOrion"},
			outMessage: map[string]string{"k1": "v1", "k2": "v2", "severity": "informational", "vendor_severity": "Informational"},
			err:        nil,
		},
		"orion Notice": {
			input:      "type: STI::CE::Message::RawEventOrion\n\nk1:v1\nk2: v2\nvendor_severity: Notice",
			outHeaders: map[string]string{"type": "STI::CE::Message::RawEventOrion"},
			outMessage: map[string]string{"k1": "v1", "k2": "v2", "severity": "low", "vendor_severity": "Notice"},
			err:        nil,
		},
		"orion Warning": {
			input:      "type: STI::CE::Message::RawEventOrion\n\nk1:v1\nk2: v2\nvendor_severity: Warning",
			outHeaders: map[string]string{"type": "STI::CE::Message::RawEventOrion"},
			outMessage: map[string]string{"k1": "v1", "k2": "v2", "severity": "medium", "vendor_severity": "Warning"},
			err:        nil,
		},
		"orion Serious": {
			input:      "type: STI::CE::Message::RawEventOrion\n\nk1:v1\nk2: v2\nvendor_severity: Serious",
			outHeaders: map[string]string{"type": "STI::CE::Message::RawEventOrion"},
			outMessage: map[string]string{"k1": "v1", "k2": "v2", "severity": "high", "vendor_severity": "Serious"},
			err:        nil,
		},
		"orion Critical": {
			input:      "type: STI::CE::Message::RawEventOrion\n\nk1:v1\nk2: v2\nvendor_severity: Critical",
			outHeaders: map[string]string{"type": "STI::CE::Message::RawEventOrion"},
			outMessage: map[string]string{"k1": "v1", "k2": "v2", "severity": "critical", "vendor_severity": "Critical"},
			err:        nil,
		},
		"orion Unmapped": {
			input:      "type: STI::CE::Message::RawEventOrion\n\nk1:v1\nk2: v2\nvendor_severity: Unmapped",
			outHeaders: map[string]string{"type": "STI::CE::Message::RawEventOrion"},
			outMessage: map[string]string{"k1": "v1", "k2": "v2", "severity": "unknown", "vendor_severity": "Unmapped"},
			err:        nil,
		},
	}

	for testName, test := range tests {
		t.Logf("TestParseDoc: %s", testName)
		headers, message, err := ParseDoc(test.input)
		assert.IsType(t, test.err, err, "ParseDoc error")
		assert.Equal(t, test.outHeaders, headers, "returned headers")
		assert.Equal(t, test.outMessage, message, "returned message")
	}
}
