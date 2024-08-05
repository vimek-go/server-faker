package plugins

import "github.com/gin-gonic/gin"

type PlgHandler interface {
	Respond(c *gin.Context)
}
