package shortener

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "single digit",
			input:    1,
			expected: "1",
		},
		{
			name:     "last numeric digit",
			input:    9,
			expected: "9",
		},
		{
			name:     "first lowercase letter",
			input:    10,
			expected: "a",
		},
		{
			name:     "last lowercase letter",
			input:    35,
			expected: "z",
		},
		{
			name:     "first uppercase letter",
			input:    36,
			expected: "A",
		},
		{
			name:     "last uppercase letter",
			input:    61,
			expected: "Z",
		},
		{
			name:     "base transition",
			input:    62,
			expected: "10",
		},
		{
			name:     "base transition plus one",
			input:    63,
			expected: "11",
		},
		{
			name:     "multiple digits",
			input:    124,
			expected: "20",
		},
		{
			name:     "bigger example",
			input:    12345,
			expected: "3d7",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := EncodeBase62(tt.input)

			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestDecodeBase62(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "zero",
			expected: 0,
			input:    "0",
		},
		{
			name:     "single digit",
			input:    "1",
			expected: 1,
		},
		{
			name:     "last numeric digit",
			input:    "9",
			expected: 9,
		},
		{
			name:     "first lowercase letter",
			input:    "a",
			expected: 10,
		},
		{
			name:     "last lowercase letter",
			input:    "z",
			expected: 35,
		},
		{
			name:     "first uppercase letter",
			input:    "A",
			expected: 36,
		},
		{
			name:     "last uppercase letter",
			input:    "Z",
			expected: 61,
		},
		{
			name:     "base transition",
			input:    "10",
			expected: 62,
		},
		{
			name:     "base transition plus one",
			input:    "11",
			expected: 63,
		},
		{
			name:     "multiple digits",
			input:    "20",
			expected: 124,
		},
		{
			name:     "bigger example",
			input:    "3d7",
			expected: 12345,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := DecodeBase62(tt.input)

			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}

func TestBase62RoundTrip(t *testing.T) {
	tests := []int64{
		0,
		1,
		10,
		62,
		100,
		1000,
		12345,
		999999999,
	}

	for _, n := range tests {
		t.Run(fmt.Sprintf("%d", n), func(t *testing.T) {
			encoded := EncodeBase62(n)

			decoded, err := DecodeBase62(encoded)

			require.NoError(t, err)
			require.Equal(t, n, decoded)
		})
	}
}
