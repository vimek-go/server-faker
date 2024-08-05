package enums

type RequestLocation string

const (
	requestBody  RequestLocation = "body"
	requestURL   RequestLocation = "url"
	requestQuery RequestLocation = "query"
)

func (rl RequestLocation) String() string {
	return string(rl)
}

func (rl RequestLocation) IsValid() bool {
	switch rl {
	case requestBody, requestURL, requestQuery:
		return true
	}
	return false
}

type requestLocation struct{}

func (requestLocation) Body() RequestLocation  { return requestBody }
func (requestLocation) Query() RequestLocation { return requestQuery }
func (requestLocation) URL() RequestLocation   { return requestURL }

var RequestLocations requestLocation
