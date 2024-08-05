package transformer

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/parser/dto"

	"github.com/pkg/errors"
)

const (
	arraySize = 3
	minNumber = 1
	maxNumber = 100
)

var ErrInvalidResponseType = fmt.Errorf("invalid response type")

type transformer struct{}

type Transformer interface {
	Transform(filePath, url, responseType string) (string, error)
}

func New() Transformer {
	return &transformer{}
}

func (t *transformer) Transform(filePath, url, responseType string) (string, error) {
	respType := enums.ResponseType(responseType)
	if !respType.IsValidForEndpointGeneration() {
		return "", errors.Wrapf(
			ErrInvalidResponseType,
			"response type %s is not valid. Valid types %s, %s",
			responseType,
			enums.ResponseTypes.Static(),
			enums.ResponseTypes.Dynamic(),
		)
	}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer jsonFile.Close()
	var input any
	err = json.NewDecoder(jsonFile).Decode(&input)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	endpoint := t.parseJSON(input, url, respType)
	bytes, err := json.MarshalIndent(endpoint, "", "  ")
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return string(bytes), nil
}

func (t *transformer) parseJSON(input any, url string, responseType enums.ResponseType) *dto.Endpoints {
	endpointURL := "fil-the-url"
	if url != "" {
		endpointURL = url
	}
	endpoint := &dto.Endpoint{
		URL:    endpointURL,
		Method: "GET",
		Response: &dto.Response{
			Status: http.StatusOK,
			Type:   "dynamic",
			Format: "json",
		},
	}

	params := make(dto.Params, 0)
	t.parseParams(input, &params, "", responseType)
	endpoint.Response.Object = params
	rval := &dto.Endpoints{Endpoints: []dto.Endpoint{*endpoint}}
	return rval
}

func (t *transformer) parseParams(
	value any,
	params *dto.Params,
	key string,
	responseType enums.ResponseType,
) *dto.Params {
	if valueM, ok := value.(map[string]any); ok {
		return t.parseMap(valueM, params, responseType)
	} else if valueA, ok := value.([]any); ok {
		*params = append(*params, t.parseArray(valueA, key, responseType))
	} else {
		param := t.prepareValueParam(value, key, responseType)
		*params = append(*params, param)
	}
	return params
}

func (t *transformer) parseArray(value []any, key string, responseType enums.ResponseType) dto.Param {
	if len(value) == 0 {
		return dto.Param{
			Key:   key,
			Array: &dto.Array{Min: 0, Max: 0, Element: nil},
		}
	}
	obj := make(dto.Params, 0)
	elem := t.parseParams(value[0], &obj, key, responseType)
	return dto.Param{
		Key:   key,
		Array: &dto.Array{Min: arraySize, Max: arraySize, Element: *elem},
	}
}

func (t *transformer) parseMap(value map[string]any, params *dto.Params, responseType enums.ResponseType) *dto.Params {
	for key, val := range value {
		if valM, ok := val.(map[string]any); ok {
			obj := make(dto.Params, 0)
			t.parseMap(valM, &obj, responseType)
			param := dto.Param{Object: obj, Key: key}
			*params = append(*params, param)
		} else if valA, ok := val.([]any); ok {
			*params = append(*params, t.parseArray(valA, key, responseType))
		} else {
			param := t.prepareValueParam(val, key, responseType)
			*params = append(*params, param)
		}
	}
	return params
}

func (t *transformer) prepareValueParam(value any, key string, valueType enums.ResponseType) dto.Param {
	switch valueType {
	case enums.ResponseTypes.Static():
		return dto.Param{Static: &dto.Static{Value: value}, Key: key}
	case enums.ResponseTypes.Dynamic():
		switch value := value.(type) {
		case string:
			return dto.Param{
				Random: &dto.Random{
					Type: enums.RandomKinds.StringAll().String(),
					Min:  minNumber,
					Max:  maxNumber,
				},
				Key: key,
			}
		case float64:
			// check if the value is an integer
			if t.isIntegral(value) {
				return dto.Param{
					Random: &dto.Random{
						Type: enums.RandomKinds.Integer().String(),
						Min:  minNumber,
						Max:  maxNumber,
					},
					Key: key,
				}
			} else {
				return dto.Param{
					Random: &dto.Random{
						Type: enums.RandomKinds.Float().String(),
						Min:  minNumber,
						Max:  maxNumber,
					},
					Key: key,
				}
			}
		case bool:
			return dto.Param{Random: &dto.Random{Type: enums.RandomKinds.Boolean().String()}, Key: key}
		default:
			fmt.Printf("value type %T is not supported for dynamic response type\n", value)
		}

	default:
	}
	return dto.Param{Static: &dto.Static{Value: value}}
}

func (t *transformer) isIntegral(val float64) bool {
	return math.Mod(val, 1) == 0
}
