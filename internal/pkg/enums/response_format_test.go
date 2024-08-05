package enums_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
)

func TestResponseFormat(t *testing.T) {
	t.Parallel()
	testEnum(t, map[StringEnum]bool{
		enums.ResponseFormat("asd"):   false,
		enums.ResponseFormats.JSON():  true,
		enums.ResponseFormats.XML():   true,
		enums.ResponseFormats.Bytes(): true,
	})
}
