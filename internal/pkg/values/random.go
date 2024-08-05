package values

import (
	"math/rand"
	"time"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/tools"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type RandomValuer struct {
	keyValue
	kind enums.RandomKind
	min  int
	max  int
}

const (
	numbers          = "0123456789"
	lowerCaseLetters = "abcdefghijklmnopqrstuvwxyz"
	upperCaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	two              = 2
)

func NewRandomValuer(key string, kind string, min, max int) (Valuer, error) {
	kindEnum := enums.RandomKind(kind)
	if !kindEnum.IsValid() {
		return nil, errors.Wrapf(ErrNotHandledKind, "requested random valuer with kind %s", kind)
	}
	return &RandomValuer{
		keyValue: keyValue{key: key},
		kind:     kindEnum,
		min:      min,
		max:      max,
	}, nil
}

func (rv *RandomValuer) Generate(_ *gin.Context) (any, error) {
	var length int
	if rv.max == rv.min {
		length = rv.min
	} else {
		length = tools.GenerateRandomInt(rv.min, rv.max)
	}

	var value any
	switch rv.kind {
	case enums.RandomKinds.StringNumeric():
		value = rv.generateRandomString(numbers, length)
	case enums.RandomKinds.StringUpperCase():
		value = rv.generateRandomString(upperCaseLetters, length)
	case enums.RandomKinds.StringLowercase():
		value = rv.generateRandomString(lowerCaseLetters, length)
	case enums.RandomKinds.StringUperCaseNumber():
		value = rv.generateRandomString(upperCaseLetters+numbers, length)
	case enums.RandomKinds.StringLowercaseNumber():
		value = rv.generateRandomString(lowerCaseLetters+numbers, length)
	case enums.RandomKinds.StringAll():
		value = rv.generateRandomString(upperCaseLetters+lowerCaseLetters+numbers, length)
	case enums.RandomKinds.Integer():
		value = tools.GenerateRandomInt(rv.min, rv.max)
	case enums.RandomKinds.Float():
		value = rv.generateRandomFloat(float64(rv.min), float64(rv.max))
	case enums.RandomKinds.Boolean():
		value = rand.Intn(two) == 1
	}

	if key := rv.keyValue.Key(); key != nil {
		return map[string]any{*key: value}, nil
	}
	return value, nil
}

func (rv *RandomValuer) generateRandomString(components string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = components[rand.Intn(len(components))]
	}
	return string(b)
}

func (rv *RandomValuer) generateRandomFloat(min, max float64) float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Float64()*(max-min)
}

func (rv *RandomValuer) Type() enums.GenerationType {
	return enums.GenerationTypes.SingleValue()
}

func (rv *RandomValuer) IsNil() bool {
	return rv == nil
}
