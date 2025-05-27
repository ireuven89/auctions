package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ToStringTest struct {
	input   interface{}
	wantErr bool
	output  string
}

func TestToString(t *testing.T) {
	tests := []ToStringTest{
		{
			input:   "testString",
			wantErr: false,
			output:  "testString",
		},
		{
			input:   []byte{'t', 'h'},
			wantErr: false,
			output:  "th",
		},
		{
			input:   12,
			wantErr: false,
			output:  "12",
		},
		{
			input:   float32(12.6),
			wantErr: false,
			output:  "12.6",
		},
		{
			input:   12.6,
			wantErr: false,
			output:  "12.6",
		},
	}

	for _, test := range tests {
		out, err := ToString(test.input)
		assert.Equal(t, err != nil, test.wantErr)
		assert.Equal(t, out, test.output)
	}
}
