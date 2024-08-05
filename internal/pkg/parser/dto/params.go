package dto

import (
	"fmt"

	"github.com/vimek-go/server-faker/internal/pkg/enums"

	"github.com/pkg/errors"
)

var ErrParamNotValid = errors.New("param is not valid")

type Params []Param

type Array struct {
	Min     int     `json:"min"     validate:"required"`
	Max     int     `json:"max"     validate:"required"`
	Element []Param `json:"element" validate:"required"`
}

type Static struct {
	Value any `json:"value"`
}

type Random struct {
	Type string `json:"type"`
	Min  int    `json:"min"`
	Max  int    `json:"max"`
}

type Mapped struct {
	From  enums.RequestLocation `json:"from"  validate:"required, oneof=body url query"`
	Param string                `json:"param" validate:"required_unless=Form body,omitempty"`
	Index *int                  `json:"index" validate:"omitempty"`
	Path  string                `json:"path"  validate:"required_if=From body,omitempty"`
	As    string                `json:"as"    validate:"oneof=number string,omitempty"`
}

type Param struct {
	Key    string  `json:"key,omitempty"    validate:"omitempty"`
	Random *Random `json:"random,omitempty" validate:"omitempty"`
	Static *Static `json:"static,omitempty" validate:"omitempty"`
	Array  *Array  `json:"array,omitempty"  validate:"omitempty"`
	Mapped *Mapped `json:"mapped,omitempty" validate:"omitempty"`
	Object Params  `json:"object,omitempty" validate:"omitempty"`
}

func (p *Param) ValueType() (enums.ValueType, error) {
	if p.Array != nil {
		return enums.ValueTypes.Array(), nil
	}
	if p.Random != nil {
		return enums.ValueTypes.Random(), nil
	}
	if p.Static != nil {
		return enums.ValueTypes.Static(), nil
	}
	if p.Mapped != nil {
		return enums.ValueTypes.Mapped(), nil
	}
	if len(p.Object) > 0 {
		return enums.ValueTypes.Object(), nil
	}
	return "", errors.Wrapf(ErrParamNotValid, "param with key %s is not valid", p.Key)
}

func (p *Param) Details() string {
	if p.Array != nil {
		return fmt.Sprintf("array with min %d and max %d, element %+v", p.Array.Min, p.Array.Max, p.Array.Element)
	}
	if p.Object != nil {
		return fmt.Sprintf("object with %+v", p.Object)
	}
	if p.Random != nil {
		return fmt.Sprintf("random with type %s, min %d and max %d", p.Random.Type, p.Random.Min, p.Random.Max)
	}
	if p.Static != nil {
		return fmt.Sprintf("static with value %+v", p.Static.Value)
	}
	if p.Mapped != nil {
		return fmt.Sprintf(
			"mapped with from %s, key %s, index %+v, path %s, as %s",
			p.Mapped.From,
			p.Mapped.Param,
			p.Mapped.Index,
			p.Mapped.Path,
			p.Mapped.As,
		)
	}
	return "unknown"
}
