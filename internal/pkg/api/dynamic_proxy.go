package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/logger"
	"github.com/vimek-go/server-faker/internal/pkg/values"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var ErrNotHandledType = errors.New("type conversion not handled")

const urlReplacementPrefix = ":"

type dynamicProxy struct {
	baseProxyHandler
	urlValuers    map[string]values.Valuer
	queryValuers  map[string]values.Valuer
	payloadValuer values.Valuer
}

func NewDynamicProxy(
	handlerMethod, handlerURL, proxyMethod, proxyURL string,
	urlValuers, queryValuers map[string]values.Valuer,
	payloadValuer values.Valuer,
	headers map[string]string,
	logger logger.Logger,
) Handler {
	return &dynamicProxy{
		baseProxyHandler: baseProxyHandler{
			handlerMethod: handlerMethod,
			handlerURL:    handlerURL,
			proxyURL:      proxyURL,
			proxyMethod:   proxyMethod,
			headers:       headers,
			logger:        logger,
		},
		queryValuers:  queryValuers,
		urlValuers:    urlValuers,
		payloadValuer: payloadValuer,
	}
}

func (dp *dynamicProxy) Respond(c *gin.Context) {
	remote, err := dp.prepareURL(c)
	if err != nil {
		dp.logger.Errorf("error preparing url %+v\n", err)
		RespondWithConversionFailure(c, err)
		return
	}

	var buf []byte
	if dp.payloadValuer != nil && !dp.payloadValuer.IsNil() {
		payloadObj, err := dp.payloadValuer.Generate(c)
		if err != nil {
			dp.logger.Errorf("error generating proxy body %+v\n", err)
			RespondWithPayloadGenerationFailure(c, err)
			return
		}
		buf, err = json.Marshal(payloadObj)
		if err != nil {
			dp.logger.Error(err)
			return
		}
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		*req = *c.Request
		req.Host = remote.Host
		req.Method = dp.proxyMethod
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = remote.Path
		req.URL.RawQuery = remote.RawQuery

		for key, value := range dp.headers {
			req.Header.Set(key, value)
		}
		if dp.proxyMethod == http.MethodPost || dp.proxyMethod == http.MethodPut {
			req.Body = io.NopCloser(bytes.NewReader(buf))
		}
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}

func (dp *dynamicProxy) prepareURL(c *gin.Context) (*url.URL, error) {
	perparedURL, err := dp.prepareBaseURL(c)
	if err != nil {
		return nil, errors.Wrapf(err, "error generating url %s", dp.proxyURL)
	}
	values := perparedURL.Query()
	dp.logger.Infof("preparing query values %s, %d", perparedURL, len(dp.queryValuers))
	for k, v := range dp.queryValuers {
		genType := v.Type()
		if genType == enums.GenerationTypes.SingleValue() {
			val, err := v.Generate(c)
			if err != nil {
				return nil, err
			}
			value, err := dp.getStringValue(val)
			if err != nil {
				return nil, errors.Wrapf(err, "error generating query params for key %s", k)
			}
			dp.logger.Debugf("adding key value:  %v", value)
			values.Add(k, value)
		}

		if genType == enums.GenerationTypes.MultiValue() {
			val, err := v.Generate(c)
			if err != nil {
				return nil, err
			}
			arrayValues, err := dp.getArrayQueryValues(val)
			if err != nil {
				return nil, errors.Wrapf(err, "error generating array query params for key %s", k)
			}
			for i := range arrayValues {
				values.Add(k, arrayValues[i])
			}
		}
	}
	perparedURL.RawQuery = values.Encode()
	dp.logger.Infof("perpared url %s, %s", perparedURL.String(), perparedURL.RawQuery)
	return perparedURL, nil
}

func (dp *dynamicProxy) getArrayQueryValues(value any) ([]string, error) {
	arrayVal, ok := value.([]any)
	if !ok {
		return nil, errors.Wrapf(ErrNotHandledType, "cannot convert value to array %v", value)
	}
	rval := make([]string, len(arrayVal))
	for i := range rval {
		stringVal, err := dp.getStringValue(arrayVal[i])
		if err != nil {
			return nil, err
		}
		rval[i] = stringVal
	}
	return rval, nil
}

func (dp *dynamicProxy) getStringValue(value any) (string, error) {
	// generated := valuer.Generate()
	switch val := value.(type) {
	case int, int64, int32:
		return fmt.Sprintf("%d", val), nil
	case float64, float32:
		return fmt.Sprintf("%f", val), nil
	case string:
		return val, nil
	case []uint8:
		return string(val), nil
	default:
		return "", errors.Wrapf(ErrNotHandledType, "not handled conversion to type %s val: %v", reflect.TypeOf(value), value)
	}
}

func (dp *dynamicProxy) prepareBaseURL(c *gin.Context) (*url.URL, error) {
	rval := dp.proxyURL
	for key, valuer := range dp.urlValuers {
		value, err := valuer.Generate(c)
		if err != nil {
			return nil, err
		}
		val, err := dp.getStringValue(value)
		if err != nil {
			return nil, err
		}
		rval = strings.Replace(rval, fmt.Sprintf("%s%s", urlReplacementPrefix, key), val, -1)
	}
	return url.Parse(rval)
}
