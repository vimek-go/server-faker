package enums_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
)

func TestValueType(t *testing.T) {
	t.Parallel()
	testEnum(t, map[StringEnum]bool{
		enums.ValueType("asd"):    false,
		enums.ValueTypes.Static(): true,
		enums.ValueTypes.Array():  true,
		enums.ValueTypes.Random(): true,
		enums.ValueTypes.Object(): true,
		enums.ValueTypes.Mapped(): true,
	})
}
