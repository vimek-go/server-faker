package enums

type ResponseFormat string

const (
	json  ResponseFormat = "json"
	xml   ResponseFormat = "xml"
	bytes ResponseFormat = "bytes"
)

func (rf ResponseFormat) String() string {
	return string(rf)
}

func (rf ResponseFormat) IsValid() bool {
	switch rf {
	case json, xml, bytes:
		return true
	}
	return false
}

type responseFormats struct{}

func (responseFormats) JSON() ResponseFormat  { return json }
func (responseFormats) XML() ResponseFormat   { return xml }
func (responseFormats) Bytes() ResponseFormat { return bytes }

var ResponseFormats responseFormats
