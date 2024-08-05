package main

// go build -buildmode=plugin -o endpoint_logger.so endpoint_logger.go

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type plgHandler struct{}

func (plg *plgHandler) Respond(c *gin.Context) {
	req := c.Request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		slog.Error("Error reading request body", err)
		c.Status(http.StatusInternalServerError)
	}

	slog.Info("Incoming HTTP request")
	slog.Info("[HTTP LOGGER]", "method", req.Method)
	slog.Info("[HTTP LOGGER]", "path", req.URL.Path)
	slog.Info("[HTTP LOGGER]", "query", req.URL.RawQuery)
	slog.Info("[HTTP LOGGER]", "headers", fmt.Sprintf("%+v", req.Header))
	slog.Info("[HTTP LOGGER]", "ip", req.RemoteAddr)
	slog.Info("[HTTP LOGGER]", "payload", string(body))

	c.Status(http.StatusOK)
}

func main() {}

// exported
var PlgHandler plgHandler
