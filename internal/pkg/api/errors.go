package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
)

const (
	emptyParamGeneratedTitle = "Provided value for param is empty"
	conversionFailedTitle    = "Conversion failed"
	payloadGenerationTitle   = "Payload generation failed"
)

type ErrorResponse struct {
	Title  string   `json:"title"`
	URL    string   `json:"url"`
	Errors []*Erorr `json:"errors"`
}

type Erorr struct {
	Details string `json:"details"`
}

func RespondWithErrorMappingParam(c *gin.Context, err error) {
	returnErrors(c, http.StatusBadRequest, c.Request.URL.EscapedPath(), emptyParamGeneratedTitle, err)
}

func RespondWithConversionFailure(c *gin.Context, err error) {
	returnErrors(c, http.StatusBadRequest, c.Request.URL.EscapedPath(), conversionFailedTitle, err)
}

func RespondWithPayloadGenerationFailure(c *gin.Context, err error) {
	returnErrors(c, http.StatusBadRequest, c.Request.URL.EscapedPath(), payloadGenerationTitle, err)
}

func returnErrors(c *gin.Context, status int, url, title string, err error) {
	var body ErrorResponse
	var merr *multierror.Error
	if errors.As(err, &merr) {
		body = NewAPIErrorResponse(url, title, merr.Errors...)
	} else {
		body = NewAPIErrorResponse(url, title, err)
	}
	c.AbortWithStatusJSON(status, body)
}

func NewAPIErrorResponse(url, title string, errs ...error) ErrorResponse {
	rval := make([]*Erorr, len(errs))
	for i := range errs {
		rval[i] = &Erorr{
			Details: errs[i].Error(),
		}
	}
	return ErrorResponse{Title: title, URL: url, Errors: rval}
}
