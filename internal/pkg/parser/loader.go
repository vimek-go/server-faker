package parser

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/vimek-go/server-faker/internal/pkg/api"
	"github.com/vimek-go/server-faker/internal/pkg/logger"
	"github.com/vimek-go/server-faker/internal/pkg/parser/dto"
	"github.com/vimek-go/server-faker/internal/pkg/tools"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

var ErrValidation = errors.New("validation failed")

type loader struct {
	validator *validator.Validate
	factory   Factory
	logger    logger.Logger
}

type Loader interface {
	LoadConfig(filePath string) ([]api.Handler, error)
}

func NewLoader(factory Factory, logger logger.Logger) Loader {
	return &loader{validator: validator.New(validator.WithRequiredStructEnabled()), factory: factory, logger: logger}
}

func (l *loader) LoadConfig(filePath string) ([]api.Handler, error) {
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, tools.LogAndReturnError(l.logger, err, "unable to open file %s", filePath)
	}
	defer configFile.Close()
	byteValue, err := io.ReadAll(configFile)
	if err != nil {
		return nil, tools.LogAndReturnError(l.logger, err, "unable to read file %s", filePath)
	}
	var endpoints dto.Endpoints

	err = json.Unmarshal(byteValue, &endpoints)
	if err != nil {
		return nil, tools.LogAndReturnError(l.logger, err, "unable to unmarshal")
	}

	baseDir := filepath.Dir(filePath)
	rval, err := l.processConfig(baseDir, endpoints)
	if err != nil {
		return nil, err
	}
	l.logger.Infof("prepared endpoints count %d", len(rval))
	return rval, nil
}

func (l *loader) processConfig(baseDir string, endpoints dto.Endpoints) ([]api.Handler, error) {
	// loop once to validate all the endpoints format
	for i := range endpoints.Endpoints {
		errs := l.validator.Struct(&endpoints.Endpoints[i])
		if errs != nil {
			l.logger.Errorf("validation errors for url %s", endpoints.Endpoints[i].URL)
			var invalidValidationError *validator.InvalidValidationError
			if errors.As(errs, &invalidValidationError) {
				l.logger.Errorf("invalid validation error")
				return nil, ErrValidation
			}

			var validationErrors validator.ValidationErrors
			if errors.As(errs, &validationErrors) {
				for _, err := range validationErrors {
					l.logger.Errorf(
						"[Validation error]: Key: '%v': failed validation for '%v' on the '%v' tag",
						err.Field(),
						err.Value(),
						err.ActualTag(),
					)
				}
			}
			return nil, ErrValidation
		}
	}

	var fileErrors *multierror.Error
	rval := make([]api.Handler, len(endpoints.Endpoints))
	for i, e := range endpoints.Endpoints {
		l.logger.Infof("processing endpoint %s %s", e.Method, e.URL)
		if apiHandler, err := l.factory.CreateEndpoint(e, baseDir); err != nil {
			fileErrors = multierror.Append(fileErrors, err)
		} else {
			rval[i] = apiHandler
		}
	}
	if fileErrors != nil {
		for _, err := range fileErrors.Errors {
			l.logger.Error(err)
		}
		return nil, fileErrors
	}
	l.logger.Infof("parsed endpoins count %d", len(rval))
	return rval, nil
}
