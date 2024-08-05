package api

import (
	"github.com/vimek-go/server-faker/internal/pkg/enums"

	"github.com/gin-gonic/gin"
)

//nolint:unused // This is a contract that is used to mock the valuer in tests
type valuer interface {
	Generate(c *gin.Context) (interface{}, error)
	Type() enums.GenerationType
	IsNil() bool
}
