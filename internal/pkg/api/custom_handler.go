package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type customHandler struct {
	HandlerMethod    string
	HandlerURL       string
	ResponseFunction func(c *gin.Context)
}

func NewCustomHandler(url, method string, function gin.HandlerFunc) ResponseHandler {
	return &customHandler{HandlerURL: url, HandlerMethod: method, ResponseFunction: function}
}

func (ch *customHandler) Method() string {
	return ch.HandlerMethod
}

func (ch *customHandler) URL() string {
	return ch.HandlerURL
}

func (ch *customHandler) Respond(c *gin.Context) {
	ch.ResponseFunction(c)
}

func (ch *customHandler) ReturnCode() int {
	return http.StatusUnprocessableEntity
}
