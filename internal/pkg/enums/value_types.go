package enums

type ValueType string

const (
	staticValue ValueType = "static"
	arrayValue  ValueType = "array"
	randomValue ValueType = "random"
	objectValue ValueType = "object"
	mappedValue ValueType = "mapped"
)

func (et ValueType) String() string {
	return string(et)
}

func (et ValueType) IsValid() bool {
	switch et {
	case staticValue, arrayValue, randomValue, objectValue, mappedValue:
		return true
	}
	return false
}

type valueTypes struct{}

func (valueTypes) Static() ValueType { return staticValue }
func (valueTypes) Array() ValueType  { return arrayValue }
func (valueTypes) Random() ValueType { return randomValue }
func (valueTypes) Object() ValueType { return objectValue }
func (valueTypes) Mapped() ValueType { return mappedValue }

var ValueTypes valueTypes
