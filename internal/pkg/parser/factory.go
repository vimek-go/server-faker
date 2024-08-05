package parser

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"

	"github.com/vimek-go/server-faker/internal/pkg/api"
	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/logger"
	"github.com/vimek-go/server-faker/internal/pkg/parser/dto"
	"github.com/vimek-go/server-faker/internal/pkg/plugins"
	"github.com/vimek-go/server-faker/internal/pkg/tools"
	"github.com/vimek-go/server-faker/internal/pkg/values"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	ErrNotHandled       = errors.New("not handled endpoint type")
	ErrEmptyKey         = errors.New("empty param key")
	ErrDuplicatedKey    = errors.New("dupliacted param key")
	ErrNotHandledValuer = errors.New("valuer type is not handled")
)

type factory struct {
	loader plugins.Loader
	logger logger.Logger
}

type Factory interface {
	CreateEndpoint(endpoint dto.Endpoint, baseDir string) (api.Handler, error)
	CreateResponseEndpoint(endpoint dto.Endpoint, baseDir string) (api.ResponseHandler, error)
	CreateProxyEndpoint(endpoint dto.Endpoint) (api.Handler, error)
}

type pluginLoader interface {
	Load(path string) (func(*gin.Context), error)
}

func NewFactory(loader pluginLoader, logger logger.Logger) Factory {
	return &factory{loader: loader, logger: logger}
}

func (f *factory) CreateEndpoint(endpoint dto.Endpoint, baseDir string) (api.Handler, error) {
	if endpoint.Proxy != nil {
		f.logger.Info("Attempting creation of proxy endpoint")
		return f.CreateProxyEndpoint(endpoint)
	}
	if endpoint.Response != nil {
		f.logger.Info("Attempting creation of response endpoint")
		return f.CreateResponseEndpoint(endpoint, baseDir)
	}
	return nil, errors.Wrapf(ErrNotHandled, "creation requested for endpoint %+v", endpoint)
}

func (f *factory) CreateResponseEndpoint(
	endpoint dto.Endpoint,
	baseDir string,
) (handler api.ResponseHandler, err error) {
	f.logger.Infof("processing endpoint type: %s, url: %s", endpoint.Response.Type, endpoint.URL)
	switch endpoint.Response.Type {
	case enums.ResponseTypes.Static():
		responseBytes, err := f.prepareStaticBytes(
			baseDir,
			endpoint.Response.File,
			endpoint.Response.Static,
			endpoint.Response.Format,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "ertor parsing endpoint %s %s to json", endpoint.Method, endpoint.URL)
		}
		return api.NewStaticHandler(
			endpoint.Response.Format,
			endpoint.Method,
			endpoint.URL,
			endpoint.Response.Status,
			responseBytes,
			endpoint.Response.ContentType,
			f.logger,
		)
	case enums.ResponseTypes.Dynamic():
		valuer, err := f.PrepareValuer(endpoint.Response.Object, endpoint.URL)
		if err != nil {
			return nil, err
		}
		return api.NewDynamicHandler(
			endpoint.Response.Format,
			endpoint.Method,
			endpoint.URL,
			endpoint.Response.Status,
			valuer,
			f.logger,
		)

	case enums.ResponseTypes.Custom():
		responseFunction, err := f.loader.Load(filepath.Join(baseDir, endpoint.Response.File))
		if err != nil {
			return nil, errors.Wrapf(err, "error loading custom response function %s", endpoint.Response.File)
		}
		return api.NewCustomHandler(endpoint.URL, endpoint.Method, responseFunction), nil
	default:
		f.logger.Info("Not implemented yet")
		return nil, errors.New("not implemented yet")
	}
}

func (f *factory) CreateProxyEndpoint(endpoint dto.Endpoint) (handler api.Handler, err error) {
	proxy := endpoint.Proxy
	switch proxy.Type {
	case enums.ResponseTypes.Static():
		handler = api.NewStaticProxy(
			endpoint.Method,
			endpoint.URL,
			proxy.Method,
			proxy.URL,
			proxy.Headers,
			f.logger,
		)
	case enums.ResponseTypes.Dynamic():
		urlValuers, err := f.prepareProxyValuersMap(
			proxy.URLParams,
			[]enums.ValueType{enums.ValueTypes.Array(), enums.ValueTypes.Object()},
			endpoint.URL,
		)
		if err != nil {
			f.logger.Error(err)
			return nil, err
		}
		queryValuers, err := f.prepareProxyValuersMap(
			proxy.Query,
			[]enums.ValueType{enums.ValueTypes.Object()},
			endpoint.URL,
		)
		if err != nil {
			f.logger.Error(err)
			return nil, err
		}
		payloadValuer, err := f.PrepareValuer(proxy.Object, endpoint.URL)
		if err != nil {
			f.logger.Error(err)
			return nil, err
		}
		handler = api.NewDynamicProxy(
			endpoint.Method,
			endpoint.URL,
			proxy.Method,
			proxy.URL,
			urlValuers,
			queryValuers,
			payloadValuer,
			proxy.Headers,
			f.logger,
		)
	}
	return handler, err
}

func (f *factory) prepareStaticBytes(
	baseDir, pathToFile string,
	object any,
	responseFormat enums.ResponseFormat,
) ([]byte, error) {
	if object != nil && responseFormat == enums.ResponseFormats.JSON() {
		bytes, err := json.Marshal(object)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	}
	return f.loadFile(baseDir, pathToFile, responseFormat)
}

func (f *factory) loadFile(baseDir, pathToFile string, responseFormat enums.ResponseFormat) ([]byte, error) {
	if pathToFile == "" {
		return []byte{}, errors.Wrapf(ErrValidation, "file path is empty")
	}
	// to do add option to check for other key
	filePath := filepath.Join(baseDir, pathToFile)
	f.logger.Infof("reading file with data %s", filePath)
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open file: %s", filePath)
	}
	defer configFile.Close()
	byteValue, err := io.ReadAll(configFile)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read file: %s", filePath)
	}
	var validateFunc func([]byte, any) error
	f.logger.Infof("preparing validation %s", responseFormat)
	switch responseFormat {
	case enums.ResponseFormats.JSON():
		validateFunc = json.Unmarshal
	case enums.ResponseFormats.XML():
		validateFunc = xml.Unmarshal
	}
	if validateFunc != nil {
		var response interface{}
		if err := validateFunc(byteValue, &response); err != nil {
			return nil, err
		}
	}
	return byteValue, nil
}

func (f *factory) PrepareValuer(params dto.Params, url string) (values.Valuer, error) {
	// if there is only 1 param
	f.logger.Debugf("preparing valuer %+v", params)
	if len(params) > 0 {
		if len(params) == 1 {
			return f.buildValuer(params[0], url)
		}
		return f.buildObjectValuer(dto.Param{Object: params}, url)
	}
	return nil, nil
}

func (f *factory) buildValuer(param dto.Param, url string) (values.Valuer, error) {
	paramType, err := param.ValueType()
	if err != nil {
		return nil, err
	}
	f.logger.Infof("creating valuer for type %s key: %s", paramType, param.Key)
	switch paramType {
	case enums.ValueTypes.Array():
		valuer, err := f.PrepareValuer(param.Array.Element, url)
		if err != nil {
			return nil, err
		}
		return values.NewArrayValuer(param.Key, param.Array.Min, param.Array.Max, valuer), nil
	case enums.ValueTypes.Object():
		return f.buildObjectValuer(param, url)
	case enums.ValueTypes.Random():
		return values.NewRandomValuer(param.Key, param.Random.Type, param.Random.Min, param.Random.Max)
	case enums.ValueTypes.Static():
		return values.NewStaticValuer(param.Key, param.Static.Value), nil
	case enums.ValueTypes.Mapped():
		f.logger.Info("========================")
		f.logger.Info("Mapped")
		f.logger.Info("url: ", url)
		f.logger.Info("========================")
		return values.NewMappedValuer(
			param.Key,
			param.Mapped.Param,
			param.Mapped.Path,
			url,
			param.Mapped.From,
			param.Mapped.Index,
			enums.NewConvertsionType(param.Mapped.As),
			f.logger,
		)
	}
	return nil, errors.New("not implemented yet")
}

func (f *factory) prepareProxyValuersMap(
	params dto.Params,
	disabledTypes []enums.ValueType,
	fakerURL string,
) (map[string]values.Valuer, error) {
	rval := make(map[string]values.Valuer, len(params))
	duplicationMap := make(map[string]bool)
	for _, p := range params {
		key := p.Key
		if _, ok := duplicationMap[key]; ok {
			return nil, errors.Wrapf(ErrDuplicatedKey, "key: %s is duplicated in query params", p.Key)
		}
		valuer, err := f.prepareProxyURLValuer(p, disabledTypes, fakerURL)
		if err != nil {
			return nil, err
		}
		rval[key] = valuer
	}
	return rval, nil
}

func (f *factory) prepareProxyURLValuer(
	param dto.Param,
	disabledTypes []enums.ValueType,
	fakerURL string,
) (values.Valuer, error) {
	if len(param.Key) == 0 {
		return nil, errors.Wrapf(ErrEmptyKey, "url params cannot have empty keys")
	}
	if valueType, err := param.ValueType(); err != nil {
		return nil, err
	} else if tools.ArrayContains(valueType, disabledTypes) {
		return nil, errors.Wrapf(ErrNotHandled, "cannot handle valuer type: %s", valueType)
	}
	// set param key to be empty to generate only value
	param.Key = ""
	return f.buildValuer(param, fakerURL)
}

func (f *factory) buildObjectValuer(param dto.Param, url string) (values.Valuer, error) {
	valuers := make([]values.Valuer, len(param.Object))
	duplicationMap := make(map[string]bool)
	for i, p := range param.Object {
		if len(p.Key) == 0 {
			return nil, errors.Wrapf(
				ErrEmptyKey,
				"endpoint url %s, cannot have empty keys details: %s",
				url,
				p.Details(),
			)
		}
		if _, ok := duplicationMap[p.Key]; ok {
			return nil, errors.Wrapf(ErrDuplicatedKey, "endpoint url %s, key: %s is duplicated", url, p.Key)
		}
		valuer, err := f.buildValuer(p, url)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"failed building valuer type: %s for key: %s",
				enums.ValueTypes.Object(),
				p.Key,
			)
		}
		valuers[i] = valuer
	}
	return values.NewObjectValuer(param.Key, valuers), nil
}
