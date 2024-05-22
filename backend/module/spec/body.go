package spec

import "github.com/apicat/apicat/v2/backend/module/spec/jsonschema"

type Body struct {
	Schema   *jsonschema.Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	Examples map[string]Example `json:"examples,omitempty" yaml:"examples,omitempty"`
}

type Example struct {
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	Value   string `json:"value,omitempty" yaml:"value,omitempty"`
}

type HTTPBody map[string]*Body

func (b *HTTPBody) SetXDiff(x string) {
	for _, v := range *b {
		v.Schema.SetXDiff(x)
	}
}
