package openapi

import "encoding/json"

const DefaultSpecURL = "https://esi.evetech.net/meta/openapi.json"

// Spec is a minimal OpenAPI 3 document subset used by generators.
type Spec struct {
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
}

type Components struct {
	Schemas    map[string]Schema    `json:"schemas"`
	Parameters map[string]Parameter `json:"parameters"`
}

type PathItem map[string]Operation

type Operation struct {
	OperationID string              `json:"operationId"`
	Summary     string              `json:"summary"`
	Parameters  []ParameterRef      `json:"parameters"`
	RequestBody *RequestBody        `json:"requestBody"`
	Responses   map[string]Response `json:"responses"`
}

type ParameterRef struct {
	Ref string `json:"$ref"`

	Name        string          `json:"name"`
	In          string          `json:"in"`
	Required    bool            `json:"required"`
	Description string          `json:"description"`
	Schema      *SchemaRef      `json:"schema"`
}

type SchemaRef struct {
	Ref string `json:"$ref"`

	Type       SchemaType           `json:"type"`
	Format     string               `json:"format"`
	Enum       []any                `json:"enum"`
	Items      *SchemaRef           `json:"items"`
	Properties map[string]SchemaRef `json:"properties"`
	Required   []string             `json:"required"`
	OneOf      []SchemaRef          `json:"oneOf"`
	Examples   json.RawMessage      `json:"examples"`
	XCommonModel Boolish            `json:"x-common-model"`
	XEnumDescriptions []string      `json:"x-enum-descriptions"`
}

type Schema struct {
	Ref string `json:"$ref"`

	Type              SchemaType            `json:"type"`
	Format            string                `json:"format"`
	Enum              []any                 `json:"enum"`
	Items             *SchemaRef            `json:"items"`
	Properties        map[string]SchemaRef  `json:"properties"`
	Required          []string              `json:"required"`
	OneOf             []SchemaRef           `json:"oneOf"`
	Examples          json.RawMessage       `json:"examples"`
	XCommonModel      Boolish               `json:"x-common-model"`
	XEnumDescriptions []string              `json:"x-enum-descriptions"`
}

type Parameter struct {
	Name        string     `json:"name"`
	In          string     `json:"in"`
	Required    bool       `json:"required"`
	Description string     `json:"description"`
	Schema      *SchemaRef `json:"schema"`
}

type RequestBody struct {
	Required bool                      `json:"required"`
	Content  map[string]MediaTypeObject `json:"content"`
}

type MediaTypeObject struct {
	Schema SchemaRef `json:"schema"`
}

type Response struct {
	Content map[string]MediaTypeObject `json:"content"`
}
