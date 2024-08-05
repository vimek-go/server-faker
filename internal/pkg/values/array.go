package values

import (
	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/tools"

	"github.com/gin-gonic/gin"
)

type ArrayValuer struct {
	keyValue
	min    int
	max    int
	valuer Valuer
}

func NewArrayValuer(key string, min, max int, valuer Valuer) Valuer {
	return &ArrayValuer{keyValue: keyValue{key: key}, min: min, max: max, valuer: valuer}
}

func (av *ArrayValuer) Generate(c *gin.Context) (any, error) {
	var length int
	if av.max == av.min {
		length = av.min
	} else {
		length = tools.GenerateRandomInt(av.min, av.max)
	}

	values := make([]any, length)
	for i := range values {
		val, err := av.valuer.Generate(c)
		if err != nil {
			return nil, err
		}
		values[i] = val
	}

	if key := av.Key(); key != nil {
		return map[string]any{*key: values}, nil
	}
	return values, nil
}

func (av *ArrayValuer) Type() enums.GenerationType {
	return enums.GenerationTypes.MultiValue()
}

func (av *ArrayValuer) IsNil() bool {
	return av.valuer.IsNil()
}
