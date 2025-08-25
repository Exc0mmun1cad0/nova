package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type encodeTest[T int | string | []string] struct {
	name  string
	input T
	want  []byte
}

func TestEncodeSimpleString(t *testing.T) {
	var tests = []encodeTest[string]{
		{
			name:  "OK",
			input: "helloworld",
			want:  []byte("+helloworld\r\n"),
		},
		{
			name:  "OK",
			input: "jkfldvfigjrso",
			want:  []byte("+jkfldvfigjrso\r\n"),
		},
		{
			name:  "Empty string",
			input: "",
			want:  []byte("+\r\n"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, EncodeSimpleString(test.input))
		})
	}
}

func TestEncodeError(t *testing.T) {
	var tests = []encodeTest[string]{
		{
			name:  "OK",
			input: "error",
			want:  []byte("-error\r\n"),
		},
		{
			name:  "OK",
			input: "jfgldfgjskdlfgl",
			want:  []byte("-jfgldfgjskdlfgl\r\n"),
		},
		{
			name:  "Empty error",
			input: "",
			want:  []byte("-\r\n"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, EncodeError(test.input))
		})
	}
}

func TestEncodeString(t *testing.T) {
	var tests = []encodeTest[string]{
		{
			name:  "OK",
			input: "helloworld",
			want:  []byte("$10\r\nhelloworld\r\n"),
		},
		{
			name:  "OK",
			input: "jfgldfgjskdlfgl",
			want:  []byte("$15\r\njfgldfgjskdlfgl\r\n"),
		},
		{
			name:  "Empty string",
			input: "",
			want:  []byte("$0\r\n\r\n"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, EncodeString(test.input))
		})
	}
}

func TestEncodeInt(t *testing.T) {
	var tests = []encodeTest[int]{
		{
			name:  "Positive integer",
			input: 129583,
			want:  []byte(":129583\r\n"),
		},
		{
			name:  "Negative integer",
			input: -4375489,
			want:  []byte(":-4375489\r\n"),
		},
		{
			name:  "Zero",
			input: 0,
			want:  []byte(":0\r\n"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, EncodeInt(test.input))
		})
	}
}

func TestEncodeArray(t *testing.T) {
	var tests = []encodeTest[[]string]{
		{
			name:  "OK",
			input: []string{"hello", "world"},
			want:  []byte("*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"),
		},
		{
			name:  "Empty",
			input: []string{},
			want:  []byte("*0\r\n"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, EncodeArray(test.input))
		})
	}
}
