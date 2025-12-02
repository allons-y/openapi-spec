// SPDX-FileCopyrightText: Copyright 2015-2025 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package spec

import (
	"encoding/json"
	"strings"

	"github.com/go-openapi/jsonpointer"
	"github.com/go-openapi/swag/jsonutils"
)

// QueryParam creates a query parameter
func QueryParam(name string) *Parameter {
	return &Parameter{ParamProps: ParamProps{Name: name, In: "query"}}
}

// HeaderParam creates a header parameter, this is always required by default
func HeaderParam(name string) *Parameter {
	return &Parameter{ParamProps: ParamProps{Name: name, In: "header", Required: true}}
}

// PathParam creates a path parameter, this is always required
func PathParam(name string) *Parameter {
	return &Parameter{ParamProps: ParamProps{Name: name, In: "path", Required: true}}
}

// CookieParam creates a cookie parameter
func CookieParam(name string) *Parameter {
	return &Parameter{ParamProps: ParamProps{Name: name, In: "cookie"}}
}

// BodyParam is deprecated in OpenAPI v3 - use RequestBody instead
// Kept for backward compatibility but should not be used
func BodyParam(name string, schema *Schema) *Parameter {
	return &Parameter{ParamProps: ParamProps{Name: name, In: "body", Schema: schema}}
}

// FormDataParam is deprecated in OpenAPI v3 - use RequestBody with appropriate media type
// Kept for backward compatibility but should not be used
func FormDataParam(name string) *Parameter {
	return &Parameter{ParamProps: ParamProps{Name: name, In: "formData"}}
}

// FileParam is deprecated in OpenAPI v3 - use RequestBody with multipart/form-data media type
// Kept for backward compatibility but should not be used
func FileParam(name string) *Parameter {
	return &Parameter{ParamProps: ParamProps{Name: name, In: "formData"},
		SimpleSchema: SimpleSchema{Type: "file"}}
}

// SimpleArrayParam creates a param for a simple array (string, int, date etc)
func SimpleArrayParam(name, tpe, fmt string) *Parameter {
	return &Parameter{ParamProps: ParamProps{Name: name},
		SimpleSchema: SimpleSchema{Type: jsonArray, CollectionFormat: "csv",
			Items: &Items{SimpleSchema: SimpleSchema{Type: tpe, Format: fmt}}}}
}

// ParamRef creates a parameter that's a json reference
func ParamRef(uri string) *Parameter {
	p := new(Parameter)
	p.Ref = MustCreateRef(uri)
	return p
}

// ParamProps describes the specific attributes of an operation parameter (OpenAPI v3)
//
// NOTE:
// - In OpenAPI v3, valid values for "in" are: query, header, path, cookie
// - body and formData are replaced by requestBody
// - Schema is now used for all parameter types (not just body)
// - AllowEmptyValue is allowed where "in" == "query"
type ParamProps struct {
	Name            string  `json:"name,omitempty"`
	In              string  `json:"in,omitempty"`
	Description     string  `json:"description,omitempty"`
	Required        bool    `json:"required,omitempty"`
	Deprecated      bool    `json:"deprecated,omitempty"`
	AllowEmptyValue bool    `json:"allowEmptyValue,omitempty"`
	Style           string  `json:"style,omitempty"`
	Explode         *bool   `json:"explode,omitempty"`
	AllowReserved   bool    `json:"allowReserved,omitempty"`
	Schema          *Schema `json:"schema,omitempty"`
	Example         any     `json:"example,omitempty"`
	Examples        map[string]Example `json:"examples,omitempty"`
	Content         map[string]MediaType `json:"content,omitempty"`
}

// Parameter a unique parameter is defined by a combination of a [name](#parameterName) and [location](#parameterIn).
//
// In OpenAPI v3, there are four possible parameter locations (in):
//   - Path - Used together with Path Templating, where the parameter value is actually part
//     of the operation's URL. This does not include the host or base path of the API. For example, in `/items/{itemId}`,
//     the path parameter is `itemId`.
//   - Query - Parameters that are appended to the URL. For example, in `/items?id=###`, the query parameter is `id`.
//   - Header - Custom headers that are expected as part of the request.
//   - Cookie - Used to pass a specific cookie value to the API.
//
// For more information: https://spec.openapis.org/oas/v3.1.0#parameter-object
type Parameter struct {
	Refable
	CommonValidations
	SimpleSchema
	VendorExtensible
	ParamProps
}

// JSONLookup look up a value by the json property name
func (p Parameter) JSONLookup(token string) (any, error) {
	if ex, ok := p.Extensions[token]; ok {
		return &ex, nil
	}
	if token == jsonRef {
		return &p.Ref, nil
	}

	r, _, err := jsonpointer.GetForToken(p.CommonValidations, token)
	if err != nil && !strings.HasPrefix(err.Error(), "object has no field") {
		return nil, err
	}
	if r != nil {
		return r, nil
	}
	r, _, err = jsonpointer.GetForToken(p.SimpleSchema, token)
	if err != nil && !strings.HasPrefix(err.Error(), "object has no field") {
		return nil, err
	}
	if r != nil {
		return r, nil
	}
	r, _, err = jsonpointer.GetForToken(p.ParamProps, token)
	return r, err
}

// WithDescription a fluent builder method for the description of the parameter
func (p *Parameter) WithDescription(description string) *Parameter {
	p.Description = description
	return p
}

// Named a fluent builder method to override the name of the parameter
func (p *Parameter) Named(name string) *Parameter {
	p.Name = name
	return p
}

// WithLocation a fluent builder method to override the location of the parameter
func (p *Parameter) WithLocation(in string) *Parameter {
	p.In = in
	return p
}

// Typed a fluent builder method for the type of the parameter value
func (p *Parameter) Typed(tpe, format string) *Parameter {
	p.Type = tpe
	p.Format = format
	return p
}

// CollectionOf a fluent builder method for an array parameter
func (p *Parameter) CollectionOf(items *Items, format string) *Parameter {
	p.Type = jsonArray
	p.Items = items
	p.CollectionFormat = format
	return p
}

// WithDefault sets the default value on this parameter
func (p *Parameter) WithDefault(defaultValue any) *Parameter {
	p.AsOptional() // with default implies optional
	p.Default = defaultValue
	return p
}

// AllowsEmptyValues flags this parameter as being ok with empty values
func (p *Parameter) AllowsEmptyValues() *Parameter {
	p.AllowEmptyValue = true
	return p
}

// NoEmptyValues flags this parameter as not liking empty values
func (p *Parameter) NoEmptyValues() *Parameter {
	p.AllowEmptyValue = false
	return p
}

// AsOptional flags this parameter as optional
func (p *Parameter) AsOptional() *Parameter {
	p.Required = false
	return p
}

// AsRequired flags this parameter as required
func (p *Parameter) AsRequired() *Parameter {
	if p.Default != nil { // with a default required makes no sense
		return p
	}
	p.Required = true
	return p
}

// WithMaxLength sets a max length value
func (p *Parameter) WithMaxLength(maximum int64) *Parameter {
	p.MaxLength = &maximum
	return p
}

// WithMinLength sets a min length value
func (p *Parameter) WithMinLength(minimum int64) *Parameter {
	p.MinLength = &minimum
	return p
}

// WithPattern sets a pattern value
func (p *Parameter) WithPattern(pattern string) *Parameter {
	p.Pattern = pattern
	return p
}

// WithMultipleOf sets a multiple of value
func (p *Parameter) WithMultipleOf(number float64) *Parameter {
	p.MultipleOf = &number
	return p
}

// WithMaximum sets a maximum number value
func (p *Parameter) WithMaximum(maximum float64, exclusive bool) *Parameter {
	p.Maximum = &maximum
	p.ExclusiveMaximum = exclusive
	return p
}

// WithMinimum sets a minimum number value
func (p *Parameter) WithMinimum(minimum float64, exclusive bool) *Parameter {
	p.Minimum = &minimum
	p.ExclusiveMinimum = exclusive
	return p
}

// WithEnum sets a the enum values (replace)
func (p *Parameter) WithEnum(values ...any) *Parameter {
	p.Enum = append([]any{}, values...)
	return p
}

// WithMaxItems sets the max items
func (p *Parameter) WithMaxItems(size int64) *Parameter {
	p.MaxItems = &size
	return p
}

// WithMinItems sets the min items
func (p *Parameter) WithMinItems(size int64) *Parameter {
	p.MinItems = &size
	return p
}

// UniqueValues dictates that this array can only have unique items
func (p *Parameter) UniqueValues() *Parameter {
	p.UniqueItems = true
	return p
}

// AllowDuplicates this array can have duplicates
func (p *Parameter) AllowDuplicates() *Parameter {
	p.UniqueItems = false
	return p
}

// WithValidations is a fluent method to set parameter validations
func (p *Parameter) WithValidations(val CommonValidations) *Parameter {
	p.SetValidations(SchemaValidations{CommonValidations: val})
	return p
}

// UnmarshalJSON hydrates this items instance with the data from JSON
func (p *Parameter) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &p.CommonValidations); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &p.Refable); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &p.SimpleSchema); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &p.VendorExtensible); err != nil {
		return err
	}
	return json.Unmarshal(data, &p.ParamProps)
}

// MarshalJSON converts this items object to JSON
func (p Parameter) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(p.CommonValidations)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(p.SimpleSchema)
	if err != nil {
		return nil, err
	}
	b3, err := json.Marshal(p.Refable)
	if err != nil {
		return nil, err
	}
	b4, err := json.Marshal(p.VendorExtensible)
	if err != nil {
		return nil, err
	}
	b5, err := json.Marshal(p.ParamProps)
	if err != nil {
		return nil, err
	}
	return jsonutils.ConcatJSON(b3, b1, b2, b4, b5), nil
}
