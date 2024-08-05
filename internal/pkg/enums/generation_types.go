package enums

type GenerationType string

const (
	singleValue  GenerationType = "single"
	multiValue   GenerationType = "multi"
	complexValue GenerationType = "complex"
)

func (gt GenerationType) String() string {
	return string(gt)
}

func (gt GenerationType) IsValid() bool {
	switch gt {
	case singleValue, multiValue, complexValue:
		return true
	}
	return false
}

type generationType struct{}

func (generationType) SingleValue() GenerationType  { return singleValue }
func (generationType) MultiValue() GenerationType   { return multiValue }
func (generationType) ComplexValue() GenerationType { return complexValue }

var GenerationTypes generationType
