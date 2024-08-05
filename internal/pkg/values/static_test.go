package values_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/values"

	"github.com/stretchr/testify/require"
)

func TestStaticValuer_Generate(t *testing.T) {
	t.Parallel()
	const testKey = "test-key"
	testCases := []struct {
		name     string
		key      string
		input    any
		expected any
	}{
		{
			name:     "string, with key",
			key:      testKey,
			input:    "test",
			expected: map[string]any{testKey: "test"},
		},
		{
			name:     "string, no key",
			input:    "test",
			expected: "test",
		},
		{
			name:     "int, with key",
			key:      testKey,
			input:    123,
			expected: map[string]any{testKey: 123},
		},
		{
			name:     "int, no key",
			input:    123,
			expected: 123,
		},
		{
			name:     "float, with key",
			key:      testKey,
			input:    123.456,
			expected: map[string]any{testKey: 123.456},
		},
		{
			name:     "float, no key",
			input:    123.456,
			expected: 123.456,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sv := values.NewStaticValuer(tc.key, tc.input)
			val, err := sv.Generate(nil)
			require.NoError(t, err)
			require.Equal(t, tc.expected, val)
		})
	}
}
