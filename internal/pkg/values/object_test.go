package values_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/values"

	"github.com/stretchr/testify/require"
)

func TestObjectValuer_Generate(t *testing.T) {
	t.Parallel()
	const testKey = "test-key"
	testCases := []struct {
		name     string
		valuers  []values.Valuer
		key      string
		expected any
	}{
		{
			name:     "empty",
			valuers:  []values.Valuer{},
			expected: map[string]any{},
		},
		{
			name:    "empty, with key",
			valuers: []values.Valuer{},
			key:     testKey,
			expected: map[string]any{
				testKey: map[string]any{},
			},
		},
		{
			name: "single static value, with key",
			valuers: []values.Valuer{
				values.NewStaticValuer("test", "value"),
			},
			key: testKey,
			expected: map[string]any{
				testKey: map[string]any{
					"test": "value",
				},
			},
		},
		{
			name: "single static value, no key",
			valuers: []values.Valuer{
				values.NewStaticValuer("test", "value"),
			},
			expected: map[string]any{
				"test": "value",
			},
		},
		{
			name: "single static value and a static array, with key",
			valuers: []values.Valuer{
				values.NewStaticValuer("test", "value"),
				values.NewArrayValuer("array", 2, 2, values.NewStaticValuer("", "aa")),
			},
			key: testKey,
			expected: map[string]any{
				testKey: map[string]any{
					"test":  "value",
					"array": []any{"aa", "aa"},
				},
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ov := values.NewObjectValuer(tc.key, tc.valuers)
			generated, err := ov.Generate(nil)
			require.NoError(t, err)
			require.Equal(t, tc.expected, generated)
		})
	}
}
