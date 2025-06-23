package k8s

import (
	"testing"
)

func TestTruncateAndCleanName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{
			name:      "No truncation needed",
			input:     "example-name",
			maxLength: 20,
			expected:  "example-name",
		},
		{
			name:      "Truncate without trailing separator",
			input:     "example-name",
			maxLength: 7,
			expected:  "example",
		},
		{
			name:      "Truncate with trailing hyphen",
			input:     "example-name-",
			maxLength: 13,
			expected:  "example-name",
		},
		{
			name:      "Truncate with trailing underscore",
			input:     "example_name_",
			maxLength: 13,
			expected:  "example_name",
		},
		{
			name:      "Truncate with trailing dot",
			input:     "example.name.",
			maxLength: 13,
			expected:  "example.name",
		},
		{
			name:      "Truncate with multiple trailing separators",
			input:     "example-name-__..",
			maxLength: 17,
			expected:  "example-name",
		},
		{
			name:      "Truncate to zero length",
			input:     "example",
			maxLength: 0,
			expected:  "",
		},
		{
			name:      "Empty input string",
			input:     "",
			maxLength: 10,
			expected:  "",
		},
		{
			name:      "Truncate with internal separators",
			input:     "ex-ample_name.test",
			maxLength: 10,
			expected:  "ex-ample_n",
		},
		{
			name:      "Truncate with trailing separator among internal ones",
			input:     "ex-ample_name.test-",
			maxLength: 19,
			expected:  "ex-ample_name.test",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TruncateAndCleanName(tc.input, tc.maxLength)
			if result != tc.expected {
				t.Errorf("TruncateAndCleanName(%q, %d) = %q; want %q", tc.input, tc.maxLength, result, tc.expected)
			}
		})
	}
}
