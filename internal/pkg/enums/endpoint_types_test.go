package enums_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
)

func TestEndpointType(t *testing.T) {
	t.Parallel()
	testEnum(t, map[StringEnum]bool{
		enums.EndpointType("asd"):    false,
		enums.EndpointTypes.Normal(): true,
		enums.EndpointTypes.Proxy():  true,
	})
}
