package spec

import (
	"errors"
	"strconv"

	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

const NODE_HTTP_REQUEST = "apicat-http-request"

type CollectionHttpRequest struct {
	Type  string            `json:"type" yaml:"type"`
	Attrs *HttpRequestAttrs `json:"attrs" yaml:"attrs"`
}

type HttpRequestAttrs struct {
	GlobalExcepts *HttpRequestGlobalExcepts `json:"globalExcepts,omitempty" yaml:"globalExcepts,omitempty"`
	Parameters    *HTTPParameters           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Content       HTTPBody                  `json:"content,omitempty" yaml:"content,omitempty"`
}

type HttpRequestGlobalExcepts struct {
	Header []int64 `json:"header" yaml:"header"`
	Cookie []int64 `json:"cookie" yaml:"cookie"`
	Query  []int64 `json:"query" yaml:"query"`
}

func init() {
	RegisterNode(&CollectionHttpRequest{
		Type: NODE_HTTP_REQUEST,
	})
}

func NewCollectionHttpRequest() *CollectionHttpRequest {
	return &CollectionHttpRequest{
		Type: NODE_HTTP_REQUEST,
		Attrs: &HttpRequestAttrs{
			GlobalExcepts: NewHttpRequestGlobalExcepts(),
			Parameters:    NewHTTPParameters(),
			Content:       make(HTTPBody),
		},
	}
}

func NewDefaultCollectionHttpRequest() *CollectionHttpRequest {
	return &CollectionHttpRequest{
		Type: NODE_HTTP_REQUEST,
		Attrs: &HttpRequestAttrs{
			GlobalExcepts: NewHttpRequestGlobalExcepts(),
			Parameters:    NewHTTPParameters(),
		},
	}
}

func NewHttpRequestGlobalExcepts() *HttpRequestGlobalExcepts {
	return &HttpRequestGlobalExcepts{
		Header: []int64{},
		Cookie: []int64{},
		Query:  []int64{},
	}
}

func (r *CollectionHttpRequest) NodeType() string {
	return r.Type
}

func (r *CollectionHttpRequest) GetGlobalExcept(in string) []int64 {
	switch in {
	case "header":
		return r.Attrs.GlobalExcepts.Header
	case "cookie":
		return r.Attrs.GlobalExcepts.Cookie
	case "query":
		return r.Attrs.GlobalExcepts.Query
	}
	return nil
}

func (r *CollectionHttpRequest) GetGlobalExcepts() *HttpRequestGlobalExcepts {
	return r.Attrs.GlobalExcepts
}

func (r *CollectionHttpRequest) GetGlobalExceptToMap() map[string][]int64 {
	return map[string][]int64{
		"header": r.Attrs.GlobalExcepts.Header,
		"cookie": r.Attrs.GlobalExcepts.Cookie,
		"query":  r.Attrs.GlobalExcepts.Query,
	}
}

func (r *CollectionHttpRequest) AddGlobalExcept(in string, id int64) {
	switch in {
	case "header":
		if len(r.Attrs.GlobalExcepts.Header) == 0 {
			r.Attrs.GlobalExcepts.Header = append(r.Attrs.GlobalExcepts.Header, id)
		} else {
			for _, v := range r.Attrs.GlobalExcepts.Header {
				if v == id {
					return
				}
			}
			r.Attrs.GlobalExcepts.Header = append(r.Attrs.GlobalExcepts.Header, id)
		}
	case "cookie":
		if len(r.Attrs.GlobalExcepts.Cookie) == 0 {
			r.Attrs.GlobalExcepts.Cookie = append(r.Attrs.GlobalExcepts.Cookie, id)
		} else {
			for _, v := range r.Attrs.GlobalExcepts.Cookie {
				if v == id {
					return
				}
			}
			r.Attrs.GlobalExcepts.Cookie = append(r.Attrs.GlobalExcepts.Cookie, id)
		}
	case "query":
		if len(r.Attrs.GlobalExcepts.Query) == 0 {
			r.Attrs.GlobalExcepts.Query = append(r.Attrs.GlobalExcepts.Query, id)
		} else {
			for _, v := range r.Attrs.GlobalExcepts.Query {
				if v == id {
					return
				}
			}
			r.Attrs.GlobalExcepts.Query = append(r.Attrs.GlobalExcepts.Query, id)
		}
	}
}

func (r *CollectionHttpRequest) DelGlobalExcept(in string, id int64) {
	switch in {
	case "header":
		if len(r.Attrs.GlobalExcepts.Header) == 0 {
			return
		}
		for i, v := range r.Attrs.GlobalExcepts.Header {
			if v == id {
				r.Attrs.GlobalExcepts.Header = append(r.Attrs.GlobalExcepts.Header[:i], r.Attrs.GlobalExcepts.Header[i+1:]...)
				return
			}
		}
	case "cookie":
		if len(r.Attrs.GlobalExcepts.Cookie) == 0 {
			return
		}
		for i, v := range r.Attrs.GlobalExcepts.Cookie {
			if v == id {
				r.Attrs.GlobalExcepts.Cookie = append(r.Attrs.GlobalExcepts.Cookie[:i], r.Attrs.GlobalExcepts.Cookie[i+1:]...)
				return
			}
		}
	case "query":
		if len(r.Attrs.GlobalExcepts.Query) == 0 {
			return
		}
		for i, v := range r.Attrs.GlobalExcepts.Query {
			if v == id {
				r.Attrs.GlobalExcepts.Query = append(r.Attrs.GlobalExcepts.Query[:i], r.Attrs.GlobalExcepts.Query[i+1:]...)
				return
			}
		}
	}
}

func (r *CollectionHttpRequest) GetRefModelIDs() []int64 {
	ids := make([]int64, 0)
	for _, v := range r.Attrs.Content {
		if v.Schema != nil {
			ids = append(ids, v.Schema.DeepGetRefID()...)
		}
	}
	if len(ids) == 0 {
		return ids
	}

	result := make([]int64, 0)
	m := make(map[int64]bool)
	for _, v := range ids {
		if _, ok := m[v]; !ok {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}

func (r *CollectionHttpRequest) DerefGlobalParameters(params *GlobalParameters) {
	if params == nil {
		return
	}

	if len(params.Query) > 0 {
		for _, p := range params.Query {
			if r.Attrs.GlobalExcepts.Exist("query", p.ID) {
				continue
			}
			r.Attrs.Parameters.Add("query", p)
		}
		r.Attrs.GlobalExcepts.Clear("query")
	}
	if len(params.Cookie) > 0 {
		for _, p := range params.Cookie {
			if r.Attrs.GlobalExcepts.Exist("cookie", p.ID) {
				continue
			}
			r.Attrs.Parameters.Add("cookie", p)
		}
		r.Attrs.GlobalExcepts.Clear("cookie")
	}
	if len(params.Header) > 0 {
		for _, p := range params.Header {
			if r.Attrs.GlobalExcepts.Exist("header", p.ID) {
				continue
			}
			r.Attrs.Parameters.Add("header", p)
		}
		r.Attrs.GlobalExcepts.Clear("header")
	}
}

func (r *CollectionHttpRequest) DerefModel(ref *DefinitionModel) error {
	if ref == nil {
		return errors.New("model is nil")
	}
	ref.Schema.ID = ref.ID

	for _, v := range r.Attrs.Content {
		if v.Schema != nil {
			refSchemas := v.Schema.DeepFindRefById(strconv.FormatInt(ref.ID, 10))
			if len(refSchemas) > 0 {
				for _, schema := range refSchemas {
					if err := schema.ReplaceRef(ref.Schema); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (r *CollectionHttpRequest) DeepDerefModel(refs DefinitionModels) error {
	if len(refs) == 0 {
		return nil
	}

	helper := jsonschema.NewDerefHelper(refs.ToJsonSchemaMap())
	return r.DeepDerefModelByHelper(helper)
}

func (r *CollectionHttpRequest) DeepDerefModelByHelper(helper *jsonschema.DerefHelper) error {
	if helper == nil {
		return errors.New("helper is nil")
	}

	for _, v := range r.Attrs.Content {
		if v.Schema != nil {
			new, err := helper.DeepDeref(v.Schema)
			if err != nil {
				return err
			}
			v.Schema = &new
		}
	}
	return nil
}

func (r *CollectionHttpRequest) ToCollectionNode() *CollectionNode {
	return &CollectionNode{
		Node: r,
	}
}

func (r *CollectionHttpRequest) DelRefModel(ref *DefinitionModel) {
	if ref == nil {
		return
	}
	ref.Schema.ID = ref.ID

	for _, v := range r.Attrs.Content {
		if v.Schema != nil {
			v.Schema.DelRef(ref.Schema)
		}
	}
}

func (g *HttpRequestGlobalExcepts) Exist(in string, id int64) bool {
	if id == 0 {
		return false
	}

	switch in {
	case "header":
		for _, v := range g.Header {
			if v == id {
				return true
			}
		}
	case "cookie":
		for _, v := range g.Cookie {
			if v == id {
				return true
			}
		}
	case "query":
		for _, v := range g.Query {
			if v == id {
				return true
			}
		}
	}
	return false
}

func (g *HttpRequestGlobalExcepts) ToMap() map[string][]int64 {
	return map[string][]int64{
		"header": g.Header,
		"cookie": g.Cookie,
		"query":  g.Query,
	}
}

func (g *HttpRequestGlobalExcepts) Clear(in string) {
	switch in {
	case "header":
		g.Header = []int64{}
	case "cookie":
		g.Cookie = []int64{}
	case "query":
		g.Query = []int64{}
	}
}
