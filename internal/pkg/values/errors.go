package values

import "errors"

var (
	ErrFailedBindingBody     = errors.New("failed binding body")
	ErrConversionFailed      = errors.New("failed converting param")
	ErrFailedLocatingElement = errors.New("failed locating element")
	ErrEmptyKey              = errors.New("url param cannot have empty key")
	ErrNotHandledType        = errors.New("valuer type is not handled")
	ErrNotHandledKind        = errors.New("random kind is not handled")
)
