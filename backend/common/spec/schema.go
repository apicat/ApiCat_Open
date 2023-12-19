package spec

import (
	"strconv"

	"github.com/apicat/apicat/backend/common/spec/jsonschema"
)

type Referencer interface {
	Ref() bool
}

type Schema struct {
	ID          int64               `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Required    bool                `json:"required,omitempty"`
	Example     any                 `json:"example,omitempty"`
	Examples    map[string]*Example `json:"examples,omitempty"`
	Schema      *jsonschema.Schema  `json:"schema,omitempty"`
	Reference   *string             `json:"$ref,omitempty"`
	XDiff       *string             `json:"x-apicat-diff,omitempty"`
}

type Example struct {
	Summary string  `json:"summary"`
	Value   any     `json:"value"`
	XDiff   *string `json:"x-apicat-diff,omitempty"`
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

func (s *Schemas) SetXDiff(x *string) {
	for _, v := range *s {
		v.SetXDiff(x)
	}
}

func (s *Schema) SetXDiff(x *string) {
	if s.Schema != nil {
		s.Schema.SetXDiff(x)
	} else {
		s.XDiff = x
	}
}

func (s *Schema) FindExample(summary string) (*Example, bool) {
	for _, v := range s.Examples {
		if v.Summary == summary {
			return v, true
		}
	}
	return nil, false
}

func (s *Schema) EqualNomal(o *Schema) (b bool) {
	b = true
	if s.Description != o.Description || s.Required != o.Required || s.Example != o.Example {
		b = false
	}
	return b && s.EqualExamples(o)
}

func (s *Schema) EqualExamples(o *Schema) (b bool) {
	b = true
	names := map[string]struct{}{}
	for _, v := range s.Examples {
		names[v.Summary] = struct{}{}
	}
	for _, v := range o.Examples {
		names[v.Summary] = struct{}{}
	}

	for k := range names {
		se, s_has := s.FindExample(k)
		oe, o_has := o.FindExample(k)
		if !s_has && o_has {
			s := "+"
			oe.XDiff = &s
			b = false
		} else if s_has && !o_has {
			s := "-"
			se.XDiff = &s
			if o.Examples == nil {
				o.Examples = make(map[string]*Example)
			}
			o.Examples[strconv.Itoa(len(o.Examples))] = se
			b = false
		} else if se.Value != oe.Value {
			s := "!"
			oe.XDiff = &s
			b = false
		}
	}
	return b
}
