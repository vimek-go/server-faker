package enums

type EndpointType string

const (
	normal EndpointType = "normal"
	proxy  EndpointType = "proxy"
)

func (et EndpointType) String() string {
	return string(et)
}

func (et EndpointType) IsValid() bool {
	switch et {
	case normal, proxy:
		return true
	}
	return false
}

type endpointTypes struct{}

func (endpointTypes) Normal() EndpointType { return normal }
func (endpointTypes) Proxy() EndpointType  { return proxy }

var EndpointTypes endpointTypes
