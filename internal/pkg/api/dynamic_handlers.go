package api

import (
	"net/http"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/logger"
	"github.com/vimek-go/server-faker/internal/pkg/values"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewDynamicHandler(
	responseFormat enums.ResponseFormat,
	method, url string,
	responseCode int,
	valuer values.Valuer,
	logger logger.Logger,
) (ResponseHandler, error) {
	baseHandler := newBaseResponseHandler(method, url, responseCode, logger)
	switch responseFormat {
	case enums.ResponseFormats.JSON():
		return &dynamicJSONHandler{baseResponseHandler: baseHandler, valuer: valuer}, nil
	default:
		return nil, errors.Wrapf(
			ErrNotSupportedFormat,
			"format %s is not supported in dynamic endpoint",
			responseFormat,
		)
	}
}

type dynamicJSONHandler struct {
	baseResponseHandler
	valuer values.Valuer
}

func (djh *dynamicJSONHandler) Respond(c *gin.Context) {
	response, err := djh.valuer.Generate(c)
	if err != nil {
		djh.Logger.Error(err)
		switch {
		case errors.Is(err, values.ErrFailedLocatingElement):
			RespondWithErrorMappingParam(c, err)
		case errors.Is(err, values.ErrConversionFailed):
			RespondWithConversionFailure(c, err)
		default:
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	c.JSON(djh.Code, response)
}
