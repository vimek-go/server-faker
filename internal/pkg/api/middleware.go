package api

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

func GinPayloadLoggerMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, err := io.ReadAll(tee)
		if err != nil {
			// nothing to do here just log
			logger.Infof("[RequestLog] url: %s\nerror: %v\n", c.Request.URL, err)
		} else {
			c.Request.Body = io.NopCloser(&buf)
			logger.Infof("[RequestLog] url: %s %s", c.Request.Method, c.Request.URL)
			logger.Infof("[RequestLog] body: %s", string(body))
			logger.Infof("[RequestLog] headers: %v", c.Request.Header)
		}
		c.Next()
	}
}

func GinStandardLoggerMiddleware() gin.HandlerFunc {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] [%s] [%s] %s %s %s %d user agent:[%s]\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Latency,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Request.UserAgent(),
		)
	})
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinResponseLogMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		statusCode := c.Writer.Status()
		logger.Infof("[ResponseLog] url: %s %s, status: %d", c.Request.Method, c.Request.URL, statusCode)
		logger.Infof("[ResponseLog] body: %s", blw.body.String())
	}
}
