package shortener

import (
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
