package values_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/values"

	"github.com/stretchr/testify/assert"
)

func TestArrayValuer_Generate(t *testing.T) {
	t.Parallel()
	const testKey = "test-key"
	testStaticValuer := values.NewStaticValuer("", "test")
	testCases := []struct {
		name     string
		key      string
		min      int
		max      int
		expected any
	}{
		{
			name:     "empty array, with key",
			min:      0,
			max:      0,
			key:      testKey,
			expected: map[string]any{testKey: []any{}},
		},
		{
			name:     "empty array, no key",
			min:      0,
			max:      0,
			expected: []any{},
		},
		{
			name:     "min and max are the same, with key",
			min:      3,
			max:      3,
			key:      testKey,
			expected: map[string]any{testKey: []any{"test", "test", "test"}},
		},
		{
			name:     "min and max are the same, no key",
			min:      3,
			max:      3,
			expected: []any{"test", "test", "test"},
		},
	}

	for i := range testCases {
		testCase := testCases[i]
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			av := values.NewArrayValuer(testCase.key, testCase.min, testCase.max, testStaticValuer)
			generated, err := av.Generate(nil)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, generated)
		})
	}
}
