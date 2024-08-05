package values

import (
	"github.com/vimek-go/server-faker/internal/pkg/enums"

	"github.com/gin-gonic/gin"
)

type Valuer interface {
	Generate(c *gin.Context) (any, error)
	Type() enums.GenerationType
	IsNil() bool
}

type keyValue struct {
	key string
}

func (k *keyValue) Key() *string {
	if len(k.key) > 0 {
		return &k.key
	}
	return nil
}
