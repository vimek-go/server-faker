package api

import (
	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	ErrNotSupportedResponseFormat = errors.New("not supported respons type")
	ErrContentTypeEmpty           = errors.New("content type empty")
	ErrNotSupportedFormat         = errors.New("not supported format")
)

type Handler interface {
	Method() string
	URL() string
	Respond(c *gin.Context)
}

type ResponseHandler interface {
	Handler
	ReturnCode() int
}

type baseResponseHandler struct {
	HandlerMethod string
	HandlerURL    string
	Code          int
	Logger        logger.Logger
}

func newBaseResponseHandler(
	method, url string,
	code int,
	logger logger.Logger,
) baseResponseHandler {
	return baseResponseHandler{
		HandlerMethod: method,
		HandlerURL:    url,
		Code:          code,
		Logger:        logger,
	}
}

func (bh *baseResponseHandler) Method() string {
	return bh.HandlerMethod
}

func (bh *baseResponseHandler) URL() string {
	return bh.HandlerURL
}

func (bh *baseResponseHandler) ReturnCode() int {
	return bh.Code
}
