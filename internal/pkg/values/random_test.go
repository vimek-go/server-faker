package values_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/values"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	t.Parallel()
	const testKey = "test-key"

	stringAsserts := func(value any, min, max int) string {
		stringValue, ok := value.(string)
		assert.True(t, ok)
		assert.True(t, len(stringValue) >= min, "expected at least %d values have %d", min, len(stringValue))
		assert.True(t, len(stringValue) <= max, "expected at most %d values have %d", max, len(stringValue))
		return stringValue
	}

	testCases := []struct {
		name        string
		kind        string
		key         string
		expectedErr error
		min         int
		max         int
		asserts     func(t *testing.T, min, max int, val any)
	}{
		{
			name:        "not handled kind",
			kind:        "not-handled",
			expectedErr: values.ErrNotHandledKind,
		},
		{
			name: "numeric string, with key",
			key:  testKey,
			kind: enums.RandomKinds.StringNumeric().String(),
			min:  2,
			max:  10,
			asserts: func(t *testing.T, _, _ int, val any) {
				value, ok := val.(map[string]any)
				assert.True(t, ok)
				assert.NotEmpty(t, value[testKey])
				stringValue := stringAsserts(value[testKey], 2, 10)
				assert.Regexp(t, `^\d+$`, stringValue)
			},
		},
		{
			name: "numeric string, no key",
			kind: enums.RandomKinds.StringNumeric().String(),
			min:  3,
			max:  7,
			asserts: func(t *testing.T, _, _ int, val any) {
				stringValue := stringAsserts(val, 3, 7)
				assert.Regexp(t, `^\d+$`, stringValue)
			},
		},
		{
			name: "upper case string, no key",
			kind: enums.RandomKinds.StringUpperCase().String(),
			min:  3,
			max:  3,
			asserts: func(t *testing.T, _, _ int, val any) {
				stringValue := stringAsserts(val, 3, 3)
				assert.Regexp(t, `^[A-Z]+$`, stringValue)
			},
		},
		{
			name: "lower case string, no key",
			kind: enums.RandomKinds.StringLowercase().String(),
			min:  3,
			max:  3,
			asserts: func(t *testing.T, _, _ int, val any) {
				stringValue := stringAsserts(val, 3, 3)
				assert.Regexp(t, `^[a-z]+$`, stringValue)
			},
		},
		{
			name: "upper case and number string, no key",
			kind: enums.RandomKinds.StringUperCaseNumber().String(),
			min:  7,
			max:  14,
			asserts: func(t *testing.T, _, _ int, val any) {
				stringValue := stringAsserts(val, 7, 14)
				assert.Regexp(t, `^[A-Z0-9]+$`, stringValue)
			},
		},
		{
			name: "lower case and number string, no key",
			kind: enums.RandomKinds.StringLowercaseNumber().String(),
			min:  7,
			max:  14,
			asserts: func(t *testing.T, _, _ int, val any) {
				stringValue := stringAsserts(val, 7, 14)
				assert.Regexp(t, `^[a-z0-9]+$`, stringValue)
			},
		},
		{
			name: "all string, no key",
			kind: enums.RandomKinds.StringAll().String(),
			min:  7,
			max:  14,
			asserts: func(t *testing.T, _, _ int, val any) {
				stringValue := stringAsserts(val, 7, 14)
				assert.Regexp(t, `^[A-Za-z0-9]+$`, stringValue)
			},
		},
		{
			name: "integer, no key",
			kind: enums.RandomKinds.Integer().String(),
			min:  3,
			max:  7,
			asserts: func(t *testing.T, min, max int, val any) {
				intValue, ok := val.(int)
				assert.True(t, ok)
				assert.GreaterOrEqual(t, intValue, min)
				assert.LessOrEqual(t, intValue, max)
			},
		},
		{
			name: "float, no key",
			kind: enums.RandomKinds.Float().String(),
			min:  3,
			max:  7,
			asserts: func(t *testing.T, min, max int, val any) {
				floatValue, ok := val.(float64)
				assert.True(t, ok)
				assert.GreaterOrEqual(t, floatValue, float64(min))
				assert.LessOrEqual(t, floatValue, float64(max))
			},
		},
		{
			name: "boolean, no key",
			kind: enums.RandomKinds.Boolean().String(),
			asserts: func(t *testing.T, _, _ int, val any) {
				_, ok := val.(bool)
				assert.True(t, ok)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rv, err := values.NewRandomValuer(tc.key, tc.kind, tc.min, tc.max)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				return
			}
			val, err := rv.Generate(nil)
			assert.NoError(t, err)
			tc.asserts(t, tc.min, tc.max, val)
		})
	}
}
