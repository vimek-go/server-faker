package api

import (
	"net/http"
	"strings"

	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var ErrWrongPositionOfWildcard = errors.New("wildcard should be at the end of the url")

type API interface {
	Run(addr ...string) (err error)
	Engine() *gin.Engine
	AddRouters(handlers ...Handler)
}

type BaseAPI struct {
	e      *gin.Engine
	logger logger.Logger
}

func NewBaseAPI(e *gin.Engine, logger logger.Logger) *BaseAPI {
	e.Use(GinStandardLoggerMiddleware())
	e.Use(GinPayloadLoggerMiddleware(logger))
	e.Use(GinResponseLogMiddleware(logger))
	be := &BaseAPI{e: e, logger: logger}
	return be
}

func (a *BaseAPI) Run(addr ...string) (err error) {
	return a.e.Run(addr...)
}

func (a *BaseAPI) Engine() *gin.Engine {
	return a.e
}

func (a *BaseAPI) AddEndpoints(handlers []Handler) {
	a.logger.Infof("adding handlers %d", len(handlers))
	for _, h := range handlers {
		a.AddRoute(h)
	}
}

func (a *BaseAPI) AddRoute(h Handler) {
	handlerFunc := a.decorateHandlerFunc(h)
	url := h.URL()
	var err error
	if strings.Contains(h.URL(), "*") {
		url, err = a.validateRemoveVildCardFromURL(h.URL())
		if err != nil {
			a.logger.Error(err)
			return
		}
	}
	a.logger.Infof("Adding new endpoint %s: url: %s\n", h.Method(), h.URL())
	switch h.Method() {
	case http.MethodGet:
		a.e.GET(url, handlerFunc)
	case http.MethodPost:
		a.e.POST(url, handlerFunc)
	case http.MethodPut:
		a.e.PUT(url, handlerFunc)
	case http.MethodDelete:
		a.e.DELETE(url, handlerFunc)
	case http.MethodPatch:
		a.e.PATCH(url, handlerFunc)
	}
}

func (a *BaseAPI) decorateHandlerFunc(h Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.Respond(ctx)
	}
}

func (a *BaseAPI) validateRemoveVildCardFromURL(url string) (string, error) {
	if url[len(url)-1] == '*' {
		return url[:len(url)-1] + "*path", nil
	}
	return url, errors.Wrapf(ErrWrongPositionOfWildcard, "url: %s", url)
}
