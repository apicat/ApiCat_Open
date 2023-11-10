package spec

import (
	"github.com/apicat/apicat/backend/module/spec/jsonschema"
)

type Referencer interface {
	Ref() bool
}

type Schema struct {
	ID          int64              `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Required    bool               `json:"required,omitempty"`
	Schema      *jsonschema.Schema `json:"schema,omitempty"`
	Reference   *string            `json:"$ref,omitempty"`
}

func (s *Schema) Ref() bool { return s.Reference != nil }

type Schemas []*Schema

func (s *Schemas) Lookup(name string) *Schema {
	if s == nil {
		return nil
	}
	for _, v := range *s {
		if name == v.Name {
			return v
		}
	}
	return nil
}

func (s *Schemas) LookupID(id int64) *Schema {
	if s == nil {
		return nil
	}
	for _, v := range *s {
		if id == v.ID {
			return v
		}
	}
	return nil
}

func (s *Schemas) Length() int {
	return len(*s)
}
