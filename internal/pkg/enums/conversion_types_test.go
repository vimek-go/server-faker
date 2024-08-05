package enums_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
)

func TestConversionType(t *testing.T) {
	t.Parallel()
	testEnum(t, map[StringEnum]bool{
		enums.ConversionTypes.None():    true,
		enums.ConversionTypes.Number():  true,
		enums.ConversionTypes.Text():    true,
		enums.NewConvertsionType("asd"): false,
		enums.NewConvertsionType(""):    true,
	})
}
