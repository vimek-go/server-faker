package enums_test

import (
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
)

func TestRandomKind(t *testing.T) {
	t.Parallel()
	testEnum(t, map[StringEnum]bool{
		enums.RandomKind("asd"):                   false,
		enums.RandomKinds.StringNumeric():         true,
		enums.RandomKinds.StringUpperCase():       true,
		enums.RandomKinds.StringLowercase():       true,
		enums.RandomKinds.StringUperCaseNumber():  true,
		enums.RandomKinds.StringLowercaseNumber(): true,
		enums.RandomKinds.StringAll():             true,
		enums.RandomKinds.Integer():               true,
		enums.RandomKinds.Float():                 true,
		enums.RandomKinds.Boolean():               true,
	})
}
