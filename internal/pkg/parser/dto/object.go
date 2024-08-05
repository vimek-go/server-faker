package dto

import (
	"github.com/vimek-go/server-faker/internal/pkg/enums"
)

//nolint:lll // This is a DTO
type Response struct {
	Status  int                `json:"status"                 validate:"required_unless=Type custom"`
	Type    enums.ResponseType `json:"type"                   validate:"required,oneof=static dynamic custom"`
	Headers map[string]string  `json:"headers"`
	File    string             `json:"file"`
	// reserved for static object
	// has priority over file, considered only if Type is static and Format is json
	Static      interface{}          `json:"static"`
	Object      Params               `json:"object"`
	Format      enums.ResponseFormat `json:"format"                 validate:"required_unless=Type custom,omitempty,oneof=json xml bytes"`
	ContentType string               `json:"content_type,omitempty" validate:"required_if=Format bytes"`
}

type Proxy struct {
	URL         string             `json:"url"          validate:"required"`
	Method      string             `json:"method"       validate:"required"`
	Type        enums.ResponseType `json:"type"         validate:"required,oneof=static dynamic"`
	Query       Params             `json:"query_params"`
	URLParams   Params             `json:"url_params"`
	ContentType string             `json:"content_type"`
	Headers     map[string]string  `json:"headers"`
	Object      Params             `json:"object"`
}

type Endpoints struct {
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	URL      string    `json:"url"      validate:"required,startswith=/"`
	Method   string    `json:"method"   validate:"required"`
	Response *Response `json:"response" validate:"required_without=Proxy,omitempty"`
	Proxy    *Proxy    `json:"proxy"    validate:"required_without=Response,omitempty"`
}
