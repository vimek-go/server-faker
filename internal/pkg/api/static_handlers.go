package api

import (
	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type staticHandler struct {
	baseResponseHandler
	byteValue   []byte
	contentType string
}

func (sh *staticHandler) Respond(c *gin.Context) {
	c.Data(sh.Code, sh.contentType, sh.byteValue)
}

func NewStaticHandler(
	responseFormat enums.ResponseFormat,
	method, url string,
	responseCode int,
	response []byte,
	contentType string,
	logger logger.Logger,
) (ResponseHandler, error) {
	baseHandler := newBaseResponseHandler(method, url, responseCode, logger)
	switch responseFormat {
	case enums.ResponseFormats.JSON():
		return &staticHandler{
			baseResponseHandler: baseHandler,
			byteValue:           response,
			contentType:         "application/json",
		}, nil
	case enums.ResponseFormats.XML():
		return &staticHandler{
			baseResponseHandler: baseHandler,
			byteValue:           response,
			contentType:         "application/xml",
		}, nil
	case enums.ResponseFormats.Bytes():
		if len(contentType) == 0 {
			return nil, errors.Wrapf(
				ErrContentTypeEmpty,
				"handler  method: %s, url: %s of type %s cannot be empty",
				baseHandler.HandlerMethod,
				baseHandler.HandlerURL,
				enums.ResponseFormats.Bytes(),
			)
		}
		return &staticHandler{baseResponseHandler: baseHandler, byteValue: response, contentType: contentType}, nil
	}
	return nil, errors.Wrapf(
		ErrNotSupportedResponseFormat,
		"request  method: %s, URL: %s. Not supported response format %s",
		baseHandler.HandlerMethod,
		baseHandler.HandlerURL,
		responseFormat,
	)
}
