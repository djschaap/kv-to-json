package parsedoc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDoc(t *testing.T) {
	tests := map[string]struct {
		input       string
		out_headers map[string]string
		out_message map[string]string
		err         error
	}{
		"empty doc": {
			input:       "",
			out_headers: map[string]string{},
			out_message: map[string]string{},
			err:         nil,
		},
		"no headers": {
			input:       "\nk1:v1\nk2: v2",
			out_headers: map[string]string{},
			out_message: map[string]string{"k1": "v1", "k2": "v2"},
			err:         nil,
		},
		"edge cases": {
			input:       "X-key:123\n\nk0\nk1:v1\nk2: v2\nk3:\n k4:no\nk5:5\nk6:",
			out_headers: map[string]string{"X-key": "123"},
			out_message: map[string]string{"k1": "v1", "k2": "v2", "k3": "", "k5": "5", "k6": ""},
			err:         nil,
		},
	}

	for test_name, test := range tests {
		t.Logf("TestParseDoc: %s", test_name)
		headers, message, err := ParseDoc(test.input)
		assert.IsType(t, test.err, err, "ParseDoc error")
		assert.Equal(t, test.out_headers, headers, "returned headers")
		assert.Equal(t, test.out_message, message, "returned message")
	}
}
