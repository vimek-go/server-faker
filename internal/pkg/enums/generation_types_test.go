package enums_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
)

func TestGenerationType(t *testing.T) {
	t.Parallel()
	testEnum(t, map[StringEnum]bool{
		enums.GenerationType("asd"):          false,
		enums.GenerationTypes.SingleValue():  true,
		enums.GenerationTypes.MultiValue():   true,
		enums.GenerationTypes.ComplexValue(): true,
	})
}
