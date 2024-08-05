package enums_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type StringEnum interface {
	IsValid() bool
	String() string
}

func testEnum(t *testing.T, tc map[StringEnum]bool) {
	t.Helper()
	for enum, expected := range tc {
		enum, expected := enum, expected
		t.Run(enum.String(), func(t *testing.T) {
			t.Parallel()
			if enum.IsValid() != expected {
				require.Equal(t, expected, enum.IsValid())
			}
		})
	}
}
