package parser

import "github.com/gin-gonic/gin"

//nolint:unused // This is a contract that is used to mock the handler in tests
type handler interface {
	Method() string
	URL() string
	Respond(c *gin.Context)
}
