package plugins

import (
	"plugin"

	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const plgHandlerName = "PlgHandler"

var ErrConversionFailed = errors.New("cannot convert to plugin intrafece")

type loader struct {
	logger logger.Logger
}

type Loader interface {
	Load(path string) (func(*gin.Context), error)
}

func NewPlugingLoader(logger logger.Logger) Loader {
	return &loader{logger: logger}
}

func (l *loader) Load(path string) (func(*gin.Context), error) {
	l.logger.Infof("loading plugin from path %s", path)
	plg, err := plugin.Open(path)
	if err != nil {
		l.logger.Errorf("error loadin plugin from path %s: %+v", path, err)
		return nil, errors.Wrapf(err, "unable to load plugin from path %s", path)
	}

	plgHandler, err := plg.Lookup(plgHandlerName)
	if err != nil {
		l.logger.Errorf("error finding proper plugin interface: %+v", err)
		return nil, errors.Wrap(err, "unable to find proper plugin interface")
	}

	var hadler PlgHandler
	hadler, ok := plgHandler.(PlgHandler)
	if !ok {
		l.logger.Error("unable to convert to handler type")
		return nil, errors.Wrap(ErrConversionFailed, "unable to convert to handler type")
	}

	return hadler.Respond, nil
}
