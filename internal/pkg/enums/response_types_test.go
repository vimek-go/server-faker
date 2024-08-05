package enums_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
)

func TestResponseType(t *testing.T) {
	t.Parallel()
	testEnum(t, map[StringEnum]bool{
		enums.ResponseType("asd"):     false,
		enums.ResponseTypes.Static():  true,
		enums.ResponseTypes.Dynamic(): true,
		enums.ResponseTypes.Custom():  true,
	})
}
