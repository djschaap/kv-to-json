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
	}

	for testName, test := range tests {
		t.Logf("TestParseDoc: %s", testName)
		headers, message, err := ParseDoc(test.input)
		assert.IsType(t, test.err, err, "ParseDoc error")
		assert.Equal(t, test.outHeaders, headers, "returned headers")
		assert.Equal(t, test.outMessage, message, "returned message")
	}
}
