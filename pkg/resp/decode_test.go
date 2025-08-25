package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	var tests = []struct {
		name  string
		input []byte
		want  []string
		err   error
	}{
		{
			name:  "OK",
			input: []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
			want:  []string{"hello", "world"},
			err:   nil,
		},
		{
			name:  "Empty",
			input: []byte("*0\r\n"),
			want:  []string{},
			err:   nil,
		},
		{
			name:  "Invalid array length",
			input: []byte("*abc\r\n$4\r\njohn\r\n$4\n\rnwick\r\n"),
			want:  []string{},
			err:   errInvalidMultibulkLength,
		},
		{
			name:  "Invalid word length number",
			input: []byte("*3\r\n$4\r\njust\r\n$4\r\na\r\n$5\r\nstring"),
			want:  []string{},
			err:   errInvalidMultibulkFormat,
		},
		{
			name:  "Invalid word length",
			input: []byte("*2\r\n$5\r\nhellooo\r\n$5\r\nowrld"),
			want:  []string{},
			err:   errInvalidMultibulkFormat,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			result, err := Decode(test.input)
			assert.Equal(t, result, test.want)
			assert.Equal(t, err, test.err)
		})
	}
}
