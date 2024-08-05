package enums

type ConversionType string

const (
	number ConversionType = "number"
	text   ConversionType = "text"
	none   ConversionType = "none"
)

func NewConvertsionType(val string) ConversionType {
	if len(val) > 0 {
		return ConversionType(val)
	}
	return none
}

func (rt ConversionType) String() string {
	return string(rt)
}

func (rt ConversionType) IsValid() bool {
	switch rt {
	case number, text, none:
		return true
	}
	return false
}

type conversionType struct{}

func (conversionType) Number() ConversionType { return number }
func (conversionType) Text() ConversionType   { return text }
func (conversionType) None() ConversionType   { return none }

var ConversionTypes conversionType
