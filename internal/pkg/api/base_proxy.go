package api

import (
	"github.com/vimek-go/server-faker/internal/pkg/logger"
)

type baseProxyHandler struct {
	handlerMethod string
	handlerURL    string
	proxyURL      string
	proxyMethod   string
	headers       map[string]string
	logger        logger.Logger
}

func (bh *baseProxyHandler) Method() string {
	return bh.handlerMethod
}

func (bh *baseProxyHandler) URL() string {
	return bh.handlerURL
}
