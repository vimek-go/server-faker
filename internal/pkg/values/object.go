package values

import (
	"fmt"

	"github.com/vimek-go/server-faker/internal/pkg/enums"

	"github.com/gin-gonic/gin"
)

type ObjectValuer struct {
	keyValue
	valuers []Valuer
}

func NewObjectValuer(key string, valuers []Valuer) Valuer {
	return &ObjectValuer{keyValue: keyValue{key: key}, valuers: valuers}
}

func (ov *ObjectValuer) Generate(c *gin.Context) (any, error) {
	values := make(map[string]any)
	for _, valuer := range ov.valuers {
		result, err := valuer.Generate(c)
		if err != nil {
			return nil, err
		}

		generated, ok := result.(map[string]any)
		fmt.Println("converted", ok)
		if ok {
			values = ov.combineMaps(values, generated)
		}
	}
	if key := ov.Key(); key != nil {
		return map[string]any{*key: values}, nil
	}
	return values, nil
}

func (ov *ObjectValuer) combineMaps(maps ...map[string]any) map[string]any {
	combined := make(map[string]any)

	for _, m := range maps {
		for k, v := range m {
			combined[k] = v
		}
	}

	return combined
}

func (ov *ObjectValuer) Type() enums.GenerationType {
	return enums.GenerationTypes.ComplexValue()
}

func (ov *ObjectValuer) IsNil() bool {
	for _, valuer := range ov.valuers {
		if !valuer.IsNil() {
			return false
		}
	}
	return true
}
