package jsonschema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Schema struct {
	// Meta Data
	Title       string `json:"title,omitempty" yaml:"title,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Default     any    `json:"default,omitempty" yaml:"default,omitempty"`
	WriteOnly   *bool  `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	ReadOnly    *bool  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	Examples    any    `json:"examples,omitempty" yaml:"examples,omitempty"`
	Deprecated  *bool  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`

	// Core
	Reference *string `json:"$ref,omitempty" yaml:"$ref,omitempty"`

	// Applicator
	AllOf                Of                       `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	AnyOf                Of                       `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	OneOf                Of                       `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	Not                  *Schema                  `json:"not,omitempty" yaml:"not,omitempty"`
	Properties           map[string]*Schema       `json:"properties,omitempty" yaml:"properties,omitempty"`
	AdditionalProperties *ValueOrBoolean[*Schema] `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Items                *ValueOrBoolean[*Schema] `json:"items,omitempty" yaml:"items,omitempty"` // 3.1 schema or bool

	// Validation
	Type             *SchemaType              `json:"type,omitempty" yaml:"type,omitempty"` // 3.1 []string 2,3.0 string
	Enum             []any                    `json:"enum,omitempty" yaml:"enum,omitempty"`
	Pattern          string                   `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MinLength        *int64                   `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength        *int64                   `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	ExclusiveMaximum *ValueOrBoolean[float64] `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"` // 3.0 bool 3.1 int
	MultipleOf       *float64                 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	ExclusiveMinimum *ValueOrBoolean[float64] `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"` // 3.0 bool 3.1 int
	Maximum          *float64                 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	Minimum          *float64                 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	MaxProperties    *int64                   `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	MinProperties    *int64                   `json:"minProperties,omitempty" yaml:"minProperties,omitempty"`
	Required         []string                 `json:"required,omitempty" yaml:"required,omitempty"`
	MaxItems         *int64                   `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	MinItems         *int64                   `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	UniqueItems      *bool                    `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`

	// Format Annotation
	Format string `json:"format,omitempty" yaml:"format,omitempty"`

	// Extension
	ID          int64    `json:"id,omitempty" yaml:"id,omitempty"`
	XOrder      []string `json:"x-apicat-orders,omitempty" yaml:"x-apicat-orders,omitempty"`
	XMock       string   `json:"x-apicat-mock,omitempty" yaml:"x-apicat-mock,omitempty"`
	XDiff       string   `json:"x-apicat-diff,omitempty" yaml:"x-apicat-diff,omitempty"`
	XFocus      bool     `json:"x-apicat-focus,omitempty" yaml:"x-apicat-focus,omitempty"`
	XSuggestion bool     `json:"x-apicat-suggestion,omitempty" yaml:"x-apicat-suggestion,omitempty"`
	Nullable    *bool    `json:"nullable,omitempty" yaml:"nullable,omitempty"`
}

var coreTypes = []string{
	"string",
	"integer",
	"number",
	"boolean",
	"object",
	"array",
	"null",
}

func NewSchema(typ string) *Schema {
	if typ == "" {
		return &Schema{
			Type: NewSchemaType(T_OBJ),
		}
	} else {
		return &Schema{
			Type: NewSchemaType(typ),
		}
	}
}

func NewSchemaFromJson(str string) (*Schema, error) {
	if str == "" {
		return nil, errors.New("empty json content")
	}
	s := &Schema{}
	if err := json.Unmarshal([]byte(str), s); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Schema) Ref() bool { return s != nil && s.Reference != nil }

func (s *Schema) DeepRef() bool {
	if s != nil && s.Ref() {
		return true
	}

	if len(s.AllOf) > 0 {
		for _, v := range s.AllOf {
			if v.DeepRef() {
				return true
			}
		}
	}
	if len(s.AnyOf) > 0 {
		for _, v := range s.AnyOf {
			if v.DeepRef() {
				return true
			}
		}
	}
	if len(s.OneOf) > 0 {
		for _, v := range s.OneOf {
			if v.DeepRef() {
				return true
			}
		}
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			if v.DeepRef() {
				return true
			}
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		return s.Items.Value().DeepRef()
	}
	return false
}

// Check if the schema refers to this id
func (s *Schema) IsRefID(id string) bool {
	if s == nil || s.Reference == nil {
		return false
	}

	i := strings.LastIndex(*s.Reference, "/")
	if i != -1 {
		if id == (*s.Reference)[i+1:] {
			return true
		}
	}
	return false
}

func (s *Schema) DeepFindRefById(id string) (refs []*Schema) {
	if s == nil {
		return
	}

	if s.IsRefID(id) {
		refs = append(refs, s)
		return
	}

	if len(s.AllOf) > 0 {
		for _, v := range s.AllOf {
			refs = append(refs, v.DeepFindRefById(id)...)
		}
	}
	if len(s.AnyOf) > 0 {
		for _, v := range s.AnyOf {
			refs = append(refs, v.DeepFindRefById(id)...)
		}
	}
	if len(s.OneOf) > 0 {
		for _, v := range s.OneOf {
			refs = append(refs, v.DeepFindRefById(id)...)
		}
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			refs = append(refs, v.DeepFindRefById(id)...)
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		refs = append(refs, s.Items.Value().DeepFindRefById(id)...)
	}
	return
}

func (s *Schema) GetRefID() (int64, error) {
	if !s.Ref() {
		return 0, errors.New("no reference")
	}

	i := strings.LastIndex(*s.Reference, "/")
	if i != -1 {
		id, _ := strconv.ParseInt((*s.Reference)[i+1:], 10, 64)
		return id, nil
	}
	return 0, errors.New("no reference")
}

func (s *Schema) DeepGetRefID() (ids []int64) {
	if s == nil {
		return
	}

	if id, err := s.GetRefID(); err == nil {
		ids = append(ids, id)
		return
	}

	if len(s.AllOf) > 0 {
		for _, v := range s.AllOf {
			ids = append(ids, v.DeepGetRefID()...)
		}
	}
	if len(s.AnyOf) > 0 {
		for _, v := range s.AnyOf {
			ids = append(ids, v.DeepGetRefID()...)
		}
	}
	if len(s.OneOf) > 0 {
		for _, v := range s.OneOf {
			ids = append(ids, v.DeepGetRefID()...)
		}
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			ids = append(ids, v.DeepGetRefID()...)
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		ids = append(ids, s.Items.Value().DeepGetRefID()...)
	}
	return
}

func (s *Schema) ReplaceRef(ref *Schema) error {
	if !s.Ref() || ref == nil {
		return errors.New("schema is not a reference or ref is nil")
	}

	if refID, err := s.GetRefID(); err != nil || refID != ref.ID {
		return errors.New("ref id does not match")
	}

	if copyValue, err := ref.Clone(); err != nil {
		return err
	} else {
		*s = *copyValue
	}
	return nil
}

func (s *Schema) DelRef(ref *Schema) {
	if s == nil || ref == nil {
		return
	}

	if s.Ref() {
		if refID, err := s.GetRefID(); err == nil && refID == ref.ID {
			s.Reference = nil
			s.Type = NewSchemaType(T_OBJ)
			return
		}
	}

	if len(s.AllOf) > 0 {
		s.AllOf.DelRef(ref)
		if len(s.AllOf) == 0 && s.Type == nil {
			s.Type = ref.Type
		}
	}
	if len(s.AnyOf) > 0 {
		s.AnyOf.DelRef(ref)
		if len(s.AnyOf) == 0 && s.Type == nil {
			s.Type = ref.Type
		}
	}
	if len(s.OneOf) > 0 {
		s.OneOf.DelRef(ref)
		if len(s.OneOf) == 0 && s.Type == nil {
			s.Type = ref.Type
		}
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			v.DelRef(ref)
		}
	}

	if s.Items != nil && !s.Items.IsBool() && s.Items.Value() != nil {
		if s.Items.Value().IsRefID(strconv.FormatInt(ref.ID, 10)) {
			s.Items.SetValue(NewSchema(ref.Type.First()))
		}
	}
}

func (s *Schema) DelXOrderByName(name string) {
	if s == nil || name == "" {
		return
	}

	if len(s.XOrder) > 0 {
		i := 0
		for i < len(s.XOrder) {
			if s.XOrder[i] == name {
				s.XOrder = append(s.XOrder[:i], s.XOrder[i+1:]...)
				return
			}
			i++
		}
	}
}

func (s *Schema) DelRequiredByName(name string) {
	if s == nil || name == "" {
		return
	}

	if len(s.Required) > 0 {
		i := 0
		for i < len(s.Required) {
			if s.Required[i] == name {
				s.Required = append(s.Required[:i], s.Required[i+1:]...)
				return
			}
			i++
		}
	}
}

func (s *Schema) CheckAllOf() bool { return s != nil && len(s.AllOf) > 0 }

func (s *Schema) DeepCheckAllOf() bool {
	if s == nil {
		return false
	}

	if s.CheckAllOf() {
		return true
	}

	if len(s.AnyOf) > 0 {
		for _, v := range s.AnyOf {
			if v.DeepCheckAllOf() {
				return true
			}
		}
	}
	if len(s.OneOf) > 0 {
		for _, v := range s.OneOf {
			if v.DeepCheckAllOf() {
				return true
			}
		}
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			if v.DeepCheckAllOf() {
				return true
			}
		}
	}

	if s.Items != nil && !s.Items.IsBool() {
		return s.Items.Value().DeepCheckAllOf()
	}
	return false
}

func (s *Schema) MergeAllOf() {
	if s == nil {
		return
	}

	if s.CheckAllOf() {
		s.AllOf = s.AllOf.Merge()
	}

	if len(s.AnyOf) > 0 {
		for _, v := range s.AnyOf {
			v.MergeAllOf()
		}
	}
	if len(s.OneOf) > 0 {
		for _, v := range s.OneOf {
			v.MergeAllOf()
		}
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			v.MergeAllOf()
		}
	}
	if s.Items != nil && !s.Items.IsBool() {
		s.Items.Value().MergeAllOf()
	}
}

func (s *Schema) ReplaceAllOf() error {
	if s == nil {
		return errors.New("schema is nil")
	}

	if s.CheckAllOf() {
		result := s.AllOf.Merge()
		if len(result) > 1 {
			return errors.New("allOf has ref schema")
		}
		if s.Type == nil {
			s.Type = NewSchemaType(T_OBJ)
		}
		if s.Properties == nil && result[0].Properties != nil {
			s.Properties = result[0].Properties
			s.XOrder = result[0].XOrder
		}
		s.AllOf = s.AllOf[:0]
	}

	if len(s.AnyOf) > 0 {
		for _, v := range s.AnyOf {
			v.ReplaceAllOf()
		}
	}
	if len(s.OneOf) > 0 {
		for _, v := range s.OneOf {
			v.ReplaceAllOf()
		}
	}

	if s.Properties != nil {
		for _, v := range s.Properties {
			v.ReplaceAllOf()
		}
	}
	if s.Items != nil && !s.Items.IsBool() {
		s.Items.Value().ReplaceAllOf()
	}

	return nil
}

func (s *Schema) Validation(raw []byte) error {
	return nil
}

func (s *Schema) Valid() error {
	if s == nil {
		return errors.New("schema is nil")
	}

	if s.Ref() {
		return nil
	}

	for _, v := range s.Type.List() {
		if !slices.Contains(coreTypes, v) {
			return fmt.Errorf("unkowan type %s", v)
		}
		switch v {
		case "array":
			return s.checkArray()
		case "object":
			return s.checkObject()
		}
	}
	return nil
}

func (s *Schema) checkObject() error {
	if s == nil {
		return errors.New("schema is nil")
	}

	if s.Ref() || s.AdditionalProperties == nil {
		return nil
	}
	proplen := 0
	if s.Properties != nil {
		proplen = len(s.Properties)
	}
	if orderlen := len(s.XOrder); proplen > 0 {
		for name, prop := range s.Properties {
			if err := prop.Valid(); err != nil {
				return err
			}
			if orderlen > 0 {
				if !slices.Contains(s.XOrder, name) {
					return fmt.Errorf("x-apicat-order does not match the properties")
				}
			}
		}
		// check required?
	}
	if s.AdditionalProperties != nil &&
		!s.AdditionalProperties.IsBool() {
		return s.AdditionalProperties.Value().Valid()
	}
	return nil
}

func (s *Schema) checkArray() error {
	if s == nil || s.Items == nil {
		return nil
	}
	if s.Items.IsBool() {
		return nil
	}
	return s.Items.Value().Valid()
}

func (s *Schema) SetXDiff(x string) {
	if s == nil || x == "" {
		return
	}

	if len(s.AllOf) > 0 {
		for _, v := range s.AllOf {
			v.SetXDiff(x)
		}
	}
	if len(s.AnyOf) > 0 {
		for _, v := range s.AnyOf {
			v.SetXDiff(x)
		}
	}
	if len(s.OneOf) > 0 {
		for _, v := range s.OneOf {
			v.SetXDiff(x)
		}
	}
	if s.Properties != nil {
		for _, v := range s.Properties {
			v.SetXDiff(x)
		}
	}
	if s.Items != nil && !s.Items.IsBool() {
		s.Items.value.SetXDiff(x)
	}
	s.XDiff = x
}

func (s *Schema) SetDefinitionModelRef(id string) {
	ref := fmt.Sprintf("#/definitions/schemas/%s", id)
	s.Reference = &ref
}

func (s *Schema) ToJson() string {
	if s == nil {
		return ""
	}

	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Schema) Clone() (*Schema, error) {
	if s == nil {
		return nil, errors.New("schema is nil")
	}
	copyValue := &Schema{}
	if bytes, err := json.Marshal(s); err != nil {
		return nil, err
	} else {
		if err := json.Unmarshal(bytes, copyValue); err != nil {
			return nil, err
		}
	}
	return copyValue, nil
}

func (s *Schema) Equal(a *Schema) bool {
	if s == nil || a == nil {
		return false
	}

	if s.Reference != nil && a.Reference != nil {
		return *s.Reference == *a.Reference
	}

	if s.Type != nil && a.Type != nil && !s.Type.Equal(a.Type) {
		return false
	}

	if s.Title != a.Title {
		return false
	}

	if s.Type != nil && s.Type.IsOneDimensional() {
		return true
	}

	// array compare
	if s.Items != nil && a.Items != nil {
		if s.Items.IsBool() != a.Items.IsBool() {
			return false
		}
		if !s.Items.IsBool() {
			if !s.Items.Value().Equal(a.Items.Value()) {
				return false
			}
		}
		return true
	} else {
		if s.Items != nil || a.Items != nil {
			return false
		}
	}

	if len(s.AnyOf) > 0 && len(a.AnyOf) > 0 {
		if len(s.AnyOf) != len(a.AnyOf) {
			return false
		}

		stypes := make(map[string]*Schema, 0)
		for i, sv := range s.AnyOf {
			stypes[sv.Type.First()] = s.AnyOf[i]
		}
		atypes := make(map[string]*Schema, 0)
		for i, sv := range s.AnyOf {
			atypes[sv.Type.First()] = s.AnyOf[i]
		}

		if len(stypes) != len(atypes) {
			return false
		}
		for k, sv := range stypes {
			if av, ok := atypes[k]; !ok {
				return false
			} else {
				if !sv.Equal(av) {
					return false
				}
			}
		}
		return true
	} else {
		if len(s.AnyOf) > 0 || len(a.AnyOf) > 0 {
			return false
		}
	}

	if len(s.OneOf) > 0 && len(a.OneOf) > 0 {
		if len(s.OneOf) != len(a.OneOf) {
			return false
		}

		stypes := make(map[string]*Schema, 0)
		for i, sv := range s.OneOf {
			stypes[sv.Type.First()] = s.OneOf[i]
		}
		atypes := make(map[string]*Schema, 0)
		for i, sv := range s.OneOf {
			atypes[sv.Type.First()] = s.OneOf[i]
		}

		if len(stypes) != len(atypes) {
			return false
		}
		for k, sv := range stypes {
			if av, ok := atypes[k]; !ok {
				return false
			} else {
				if !sv.Equal(av) {
					return false
				}
			}
		}
		return true
	} else {
		if len(s.OneOf) > 0 || len(a.OneOf) > 0 {
			return false
		}
	}

	sProperties := make(map[string]*Schema, 0)
	if len(s.Properties) > 0 {
		for k, sv := range s.Properties {
			sProperties[k] = sv
		}
	}
	if len(s.AllOf) > 0 {
		for _, sv := range s.AllOf {
			if sv.Properties != nil {
				for k, v := range sv.Properties {
					sProperties[k] = v
				}
			}
			if refID, err := sv.GetRefID(); err == nil && refID > 0 {
				sProperties[strconv.FormatInt(refID, 10)] = sv
			}
		}
	}

	aProperties := make(map[string]*Schema, 0)
	if len(a.Properties) > 0 {
		for k, av := range a.Properties {
			aProperties[k] = av
		}
	}
	if len(a.AllOf) > 0 {
		for _, av := range a.AllOf {
			if av.Properties != nil {
				for k, v := range av.Properties {
					aProperties[k] = v
				}
			}
			if refID, err := av.GetRefID(); err == nil && refID > 0 {
				aProperties[strconv.FormatInt(refID, 10)] = av
			}
		}
	}

	if len(sProperties) != len(aProperties) {
		return false
	}
	for k, sv := range sProperties {
		if av, ok := aProperties[k]; !ok {
			return false
		} else {
			if !sv.Equal(av) {
				return false
			}
		}
	}

	return true
}
