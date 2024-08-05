package tools

import (
	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/pkg/errors"
)

func LogAndReturnError(logger logger.Logger, err error, format string, args ...any) error {
	err = errors.Wrapf(err, format, args...)
	logger.Errorf(err.Error())
	return err
}
