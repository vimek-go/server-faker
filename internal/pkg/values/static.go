package values

import (
	"github.com/vimek-go/server-faker/internal/pkg/enums"

	"github.com/gin-gonic/gin"
)

type StaticValuer[T any] struct {
	keyValue
	value T
}

func NewStaticValuer[T any](key string, value T) Valuer {
	v := StaticValuer[T]{
		keyValue: keyValue{key: key},
		value:    value,
	}
	return &v
}

func (sv *StaticValuer[T]) Generate(_ *gin.Context) (any, error) {
	if key := sv.keyValue.Key(); key != nil {
		return map[string]any{*key: sv.value}, nil
	}
	return sv.value, nil
}

func (sv *StaticValuer[T]) Type() enums.GenerationType {
	return enums.GenerationTypes.SingleValue()
}

func (sv *StaticValuer[T]) IsNil() bool {
	return sv == nil
}
