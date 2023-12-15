package diff

import (
	"fmt"

	"github.com/apicat/apicat/backend/common/spec"
	"github.com/apicat/apicat/backend/common/spec/jsonschema"
	"golang.org/x/exp/slices"
)

var (
	diffNew    = "+"
	diffRemove = "-"
	diffUpdate = "!"
	// diffType        = "type"
	// diffName        = "name"
	// diffDescription = "description"
	// diffRequired    = "required"
	// diffMock        = "mock"
	// diffDefault     = "default"
)

// Diff 比较两个接口的差异
// source,target 是完整的spec对象 因为需要解析schema等依赖
// spec.Collections 里面只能有一个接口
// 返回对比后的两个接口 其中只有最新的那个 也就是target里边会通过x-apicat-diff标记是否有差异
// 差异并不包含排序
func Diff(source, target *spec.Spec, del bool) (*spec.CollectItem, *spec.CollectItem) {
	if len(source.Collections) != 1 || len(target.Collections) != 1 {
		panic("source,target Collections length error")
	}
	a, au := getMapOne(source.CollectionsMap(true, 1))
	b, bu := getMapOne(target.CollectionsMap(true, 1))
	if au.Path != bu.Path {
		bu.XDiff = &diffUpdate
	}
	equalRequest(&a.HTTPRequestNode, &b.HTTPRequestNode, del)
	b.Responses = equalResponse(a.Responses, b.Responses, del)
	return a.ToCollectItem(*au), b.ToCollectItem(*bu)
}

func getMapOne(d map[string]map[string]spec.HTTPPart) (*spec.HTTPPart, *spec.HTTPURLNode) {
	for path, v := range d {
		for method, vv := range v {
			return &vv, &spec.HTTPURLNode{
				Method: method,
				Path:   path,
			}
		}
	}
	return nil, nil
}

func equalParam(a spec.HTTPParameters, b *spec.HTTPParameters, del bool) {
	a1 := a.Map()
	b1 := b.Map()
	for k, v := range b1 {
		x, ok := a1[k]
		if !ok {
			x = make(spec.Schemas, 0)
		}
		newv := equalSchemas(x, v, del)
		switch k {
		case "path":
			b.Path = newv
		case "header":
			b.Header = newv
		case "query":
			b.Query = newv
		case "cookie":
			b.Cookie = newv
		}
	}
}

func equalSchemas(a, b spec.Schemas, del bool) spec.Schemas {
	if del {
		for i, v := range a {
			if s := b.Lookup(v.Name); s == nil {
				newv := *v
				newv.XDiff = &diffRemove
				if i < len(b)-1 {
					b = slices.Insert(b, i, &newv)
				} else {
					b = append(b, &newv)
				}
			}
		}
	}

	for _, v := range b {
		if v.XDiff == &diffRemove {
			continue
		}
		if s := a.Lookup(v.Name); s == nil {
			v.XDiff = &diffNew
		} else {
			equalSchema(s, v, del)
		}
	}
	return b
}

func equalContent(a, b spec.HTTPBody, del bool) spec.HTTPBody {
	if del {
		for k, v := range a {
			if _, ok := b[k]; !ok {
				newv := *v
				newv.XDiff = &diffRemove
				b[k] = &newv
			}
		}
	}
	for k, v := range b {
		if x, ok := a[k]; !ok {
			v.XDiff = &diffNew
		} else {
			equalSchema(x, v, del)
		}
	}
	return b
}

func equalRequest(a, b *spec.HTTPRequestNode, del bool) {
	equalParam(a.Parameters, &b.Parameters, del)
	b.Content = equalContent(a.Content, b.Content, del)
}

func equalResponse(a, b spec.HTTPResponses, del bool) spec.HTTPResponses {
	if del {
		bb := b.Map()
		for i, v := range a {
			if _, ok := bb[v.Code]; !ok {
				newv := v
				newv.XDiff = &diffRemove
				// 尽量保证顺序
				if i < len(b)-1 {
					b = slices.Insert(b, i, newv)
				} else {
					b = append(b, newv)
				}
			}
		}
	}
	aa := a.Map()
	for i, v := range b {
		if v.XDiff == &diffRemove {
			continue
		}
		if x, ok := aa[v.Code]; ok {
			switch {
			case x.Name != v.Name || x.Description != v.Description:
				v.XDiff = &diffUpdate
			default:
				v.Header = equalSchemas(x.Header, v.Header, del)
				v.Content = equalContent(x.Content, v.Content, del)
			}
		} else {
			v.XDiff = &diffNew
		}
		b[i] = v
	}
	return b
}

func equalSchema(a, b *spec.Schema, del bool) {
	switch {
	case a.Name != b.Name:
		fmt.Println("----->>>>>")
	case a.Description != b.Description:
	case a.Required != b.Required:
	default:
		equalJsonSchema(a.Schema, b.Schema, del)
		return
	}
	b.XDiff = &diffUpdate
}

func equalJsonSchema(a, b *jsonschema.Schema, del bool) {
	if !slices.Equal(a.Type.Value(), b.Type.Value()) {
		b.XDiff = &diffUpdate
		return
	}
	at := a.Type.Value()[0]
	bt := b.Type.Value()[0]
	if at != bt {
		b.XDiff = &diffUpdate
		return
	}
	equalJsonSchemaNormal(a, b)
	switch bt {
	case "object":
		for k, v := range b.Properties {
			if x, ok := a.Properties[k]; !ok {
				v.XDiff = &diffNew
			} else {
				if slices.Contains(b.Required, k) != slices.Contains(a.Required, k) {
					v.XDiff = &diffUpdate
				} else {
					equalJsonSchema(x, v, del)
				}
			}
		}
		if del {
			for k, v := range a.Properties {
				if _, ok := b.Properties[k]; !ok {
					newv := *v
					newv.XDiff = &diffRemove
					b.Properties[k] = &newv
				}
			}
		}
	case "array":
		equalJsonSchema(a.Items.Value(), b.Items.Value(), del)
	}

}

func equalJsonSchemaNormal(a, b *jsonschema.Schema) bool {
	switch {
	case a.Default != b.Default:
	case a.Description != b.Description:
	case a.XMock != b.XMock:
	// case a.Format != b.Format:
	// case a.Pattern != b.Pattern
	default:
		return true
	}
	b.XDiff = &diffUpdate
	return false
}
