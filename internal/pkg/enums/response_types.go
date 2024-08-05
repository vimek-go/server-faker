package enums

type ResponseType string

const (
	static  ResponseType = "static"
	dynamic ResponseType = "dynamic"
	custom  ResponseType = "custom"
)

func (rt ResponseType) String() string {
	return string(rt)
}

func (rt ResponseType) IsValid() bool {
	switch rt {
	case static, dynamic, custom:
		return true
	}
	return false
}

//nolint:exhaustive // we don't need to check for all cases here
func (rt ResponseType) IsValidForEndpointGeneration() bool {
	switch rt {
	case static, dynamic:
		return true
	default:
		return false
	}
}

type responseTypes struct{}

func (responseTypes) Static() ResponseType  { return static }
func (responseTypes) Dynamic() ResponseType { return dynamic }
func (responseTypes) Custom() ResponseType  { return custom }

var ResponseTypes responseTypes
