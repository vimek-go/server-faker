package enums_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
)

func TestRequestLocation(t *testing.T) {
	t.Parallel()
	testEnum(t, map[StringEnum]bool{
		enums.RequestLocation("asd"):   false,
		enums.RequestLocations.Body():  true,
		enums.RequestLocations.Query(): true,
		enums.RequestLocations.URL():   true,
	})
}
