package values

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/PaesslerAG/jsonpath"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type PayloadMapper struct {
	keyValue
	path       string
	conversion enums.ConversionType
	logger     logger.Logger
}

func NewMappedValuer(
	responseKey, valueKey, path, url string,
	location enums.RequestLocation,
	index *int,
	conversion enums.ConversionType,
	logger logger.Logger,
) (Valuer, error) {
	switch location {
	case enums.RequestLocations.Body():
		return newPayloadMapper(responseKey, path, conversion, logger), nil
	case enums.RequestLocations.Query():
		return newQueryMapper(responseKey, valueKey, index, conversion, logger), nil
	case enums.RequestLocations.URL():
		return newURLMapper(responseKey, valueKey, url, conversion, logger)
	}
	return nil, errors.Wrapf(ErrNotHandledType, "requested valuer mapped with location %s", location)
}

func newPayloadMapper(responseKey, path string, conversion enums.ConversionType, logger logger.Logger) Valuer {
	return &PayloadMapper{keyValue: keyValue{key: responseKey}, path: path, conversion: conversion, logger: logger}
}

func (pm *PayloadMapper) Generate(c *gin.Context) (any, error) {
	body, err := getPayload(c)
	if err != nil {
		return nil, err
	}

	value, err := jsonpath.Get(pm.path, body)
	if err != nil {
		err = errors.Wrapf(ErrFailedLocatingElement, "failed getting value for json path %v", err.Error())
		pm.logger.Error(err)
		return nil, err
	}
	if pm.conversion != enums.ConversionTypes.None() {
		value, err = transform(value, pm.conversion)
		if err != nil {
			return nil, errors.Wrapf(err, "param key: [%s]", pm.key)
		}
	}
	if key := pm.keyValue.Key(); key != nil {
		return map[string]any{*key: value}, nil
	}
	return value, nil
}

func (pm *PayloadMapper) Type() enums.GenerationType {
	return enums.GenerationTypes.SingleValue()
}

func (pm *PayloadMapper) IsNil() bool {
	return pm == nil
}

type QueryMapper struct {
	keyValue
	key        string
	index      *int
	conversion enums.ConversionType
	logger     logger.Logger
}

func newQueryMapper(
	responseKey, valueKey string,
	index *int,
	conversion enums.ConversionType,
	logger logger.Logger,
) Valuer {
	return &QueryMapper{
		keyValue:   keyValue{key: responseKey},
		key:        valueKey,
		index:      index,
		conversion: conversion,
		logger:     logger,
	}
}

func (qm *QueryMapper) Generate(c *gin.Context) (any, error) {
	paramValues := c.QueryArray(qm.key)
	qm.logger.Infof("query key[%s] value %+v", qm.key, paramValues)

	if len(paramValues) == 0 {
		return nil, errors.Wrapf(ErrFailedLocatingElement, "no query argument provided for key: %s", qm.key)
	}

	paramValue := paramValues[0]
	if qm.index != nil {
		if *qm.index < 0 || *qm.index >= len(paramValues) {
			return nil, errors.Wrapf(
				ErrFailedLocatingElement,
				"index %d out of range for key: %s, array %v",
				*qm.index,
				qm.key,
				paramValues,
			)
		}
		paramValue = paramValues[*qm.index]
	}

	var value any
	var err error
	if qm.conversion != enums.ConversionTypes.None() {
		value, err = transform(paramValue, qm.conversion)
		if err != nil {
			return nil, errors.Wrapf(err, "param key: [%s]", qm.key)
		}
	} else {
		value = paramValue
	}
	if key := qm.keyValue.Key(); key != nil {
		return map[string]any{*key: value}, nil
	}
	return value, nil
}

func (qm *QueryMapper) Type() enums.GenerationType {
	return enums.GenerationTypes.SingleValue()
}

func (qm *QueryMapper) IsNil() bool {
	return qm == nil
}

type URLMapper struct {
	keyValue
	position   int
	conversion enums.ConversionType
	logger     logger.Logger
}

func newURLMapper(responseKey, key, url string, conversion enums.ConversionType, logger logger.Logger) (Valuer, error) {
	if len(key) == 0 {
		return nil, errors.Wrapf(ErrEmptyKey, "key cannot be empty for url: %s, response key: %s", url, responseKey)
	}
	position, err := findKeyInURL(key, url)
	if err != nil {
		return nil, errors.Wrapf(ErrFailedLocatingElement, "cannot locate key :%s in url %s", key, url)
	}
	return &URLMapper{
		keyValue:   keyValue{key: responseKey},
		position:   position,
		conversion: conversion,
		logger:     logger,
	}, nil
}

func (um *URLMapper) Generate(c *gin.Context) (any, error) {
	absURL := c.Request.URL.EscapedPath()
	segmentValue := um.getSegmentValueFromURL(um.position, absURL)
	var value any
	if um.conversion != enums.ConversionTypes.None() {
		var err error
		value, err = transform(segmentValue, um.conversion)
		if err != nil {
			return nil, errors.Wrapf(err, "param key: [%s]", um.key)
		}
	} else {
		value = segmentValue
	}

	if key := um.keyValue.Key(); key != nil {
		return map[string]any{*key: value}, nil
	}
	return value, nil
}

func (um *URLMapper) Type() enums.GenerationType {
	return enums.GenerationTypes.SingleValue()
}

func (um *URLMapper) getSegmentValueFromURL(segmentNumber int, url string) string {
	fragments := strings.Split(url, "/")
	// it should be safe to get element like this from url
	// in other case it should not be mached to this router
	return fragments[segmentNumber]
}

func getPayload(c *gin.Context) (any, error) {
	var body any
	if err := c.ShouldBindJSON(&body); err != nil {
		return nil, ErrFailedBindingBody
	}
	return body, nil
}

func transform(val any, conversion enums.ConversionType) (any, error) {
	switch conversion {
	case enums.ConversionTypes.Text():
		return convertToText(val)
	case enums.ConversionTypes.Number():
		return convertToNumber(val)
	}
	return val, nil
}

func convertToText(val any) (string, error) {
	switch value := val.(type) {
	// the incoming value could be of type flat64 and string
	case float64:
		return fmt.Sprintf("%g", val), nil
	case string:
		return value, nil
	}
	return "", errors.Wrapf(ErrConversionFailed, "not kwnow conversion for '%v' to string", val)
}

func convertToNumber(val any) (float64, error) {
	switch value := val.(type) {
	// the incoming value could be of type flat64 and string
	case float64:
		return value, nil
	case string:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, errors.Wrapf(ErrConversionFailed, "not known conversion for '%v' to number %v", val, err)
		}
		return f, nil
	}
	return 0, errors.Wrapf(ErrConversionFailed, "not known conversion for '%v' to number", val)
}

func findKeyInURL(key, URL string) (int, error) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return 0, err
	}

	fragments := strings.Split(parsedURL.EscapedPath(), "/")
	for i, f := range fragments {
		if strings.HasPrefix(f, ":") && strings.TrimPrefix(f, ":") == key {
			return i, nil
		}
	}
	return 0, errors.New("not found")
}

func (um *URLMapper) IsNil() bool {
	return um == nil
}
