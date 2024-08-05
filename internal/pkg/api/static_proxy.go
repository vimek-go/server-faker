package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

type staticProxy struct {
	baseProxyHandler
}

func NewStaticProxy(
	handlerMethod, handlerURL, proxyMethod, proxyURL string,
	headers map[string]string,
	logger logger.Logger,
) Handler {
	return &staticProxy{
		baseProxyHandler: baseProxyHandler{
			handlerMethod: handlerMethod,
			handlerURL:    handlerURL,
			proxyURL:      proxyURL,
			proxyMethod:   proxyMethod,
			headers:       headers,
			logger:        logger,
		},
	}
}

func (sp *staticProxy) Respond(c *gin.Context) {
	remote, err := url.Parse(sp.proxyURL)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		*req = *c.Request
		req.Host = remote.Host
		req.Method = sp.proxyMethod
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
		req.URL.RawQuery = remote.RawQuery

		for key, value := range sp.headers {
			req.Header.Set(key, value)
		}
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
