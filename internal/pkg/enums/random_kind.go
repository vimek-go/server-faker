package enums

type RandomKind string

const (
	stringNumeric   RandomKind = "string-numeric"
	upperCase       RandomKind = "string-uppercase"
	lowercase       RandomKind = "string-lowercase"
	uperCaseNumber  RandomKind = "string-uppercase-number"
	lowercaseNumber RandomKind = "string-lowercase-number"
	stringAll       RandomKind = "string-all"
	integer         RandomKind = "integer"
	float           RandomKind = "float"
	boolean         RandomKind = "boolean"
)

func (rk RandomKind) String() string {
	return string(rk)
}

func (rk RandomKind) IsValid() bool {
	switch rk {
	case stringNumeric, upperCase, lowercase, uperCaseNumber, lowercaseNumber, stringAll, integer, float, boolean:
		return true
	}
	return false
}

type randomKind struct{}

func (randomKind) StringNumeric() RandomKind         { return stringNumeric }
func (randomKind) StringUpperCase() RandomKind       { return upperCase }
func (randomKind) StringLowercase() RandomKind       { return lowercase }
func (randomKind) StringUperCaseNumber() RandomKind  { return uperCaseNumber }
func (randomKind) StringLowercaseNumber() RandomKind { return lowercaseNumber }
func (randomKind) StringAll() RandomKind             { return stringAll }
func (randomKind) Integer() RandomKind               { return integer }
func (randomKind) Float() RandomKind                 { return float }
func (randomKind) Boolean() RandomKind               { return boolean }

var RandomKinds randomKind
