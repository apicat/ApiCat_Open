package postman

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

func Import(data []byte) (*spec.Spec, error) {
	var pm Spec
	if err := json.Unmarshal(data, &pm); err != nil {
		return nil, err
	}

	p := &spec.Spec{
		ApiCat: "2.0.1",
		Info: spec.Info{
			Title:       pm.Info.Name,
			Description: pm.Info.Description,
		},
		Servers: func() []spec.Server {
			for _, v := range pm.Items {
				if v.Request != nil {
					return []spec.Server{{
						URL: fmt.Sprintf("%s://%s",
							v.Request.Url.Protocol,
							strings.Join(v.Request.Url.Host, "."),
						),
						Description: "default",
					}}
				}
			}
			return []spec.Server{}
		}(),
		Globals: func() *spec.Globals {
			return &spec.Globals{
				Parameters: spec.NewGlobalParameters(),
			}
		}(),
		Definitions: &spec.Definitions{
			Schemas:   make(spec.DefinitionModels, 0),
			Responses: make(spec.DefinitionResponses, 0),
		},
		Collections: walkCollection(pm.Items, 1000),
	}
	return p, nil
}

func walkCollection(items []Item, parentid int64) []*spec.Collection {
	cs := make(spec.Collections, 0)
	for i, v := range items {
		// http request
		id := parentid*1024 + int64(i) + 1
		if v.Request != nil {
			specItem := &spec.Collection{
				ID:       id,
				ParentID: parentid,
				Type:     spec.TYPE_HTTP,
				Title:    v.Name,
				Content:  convertContent(v),
			}
			cs = append(cs, specItem)
		}
		// dir
		if len(v.Items) > 0 {
			specItem := &spec.Collection{
				ID:       id,
				ParentID: parentid,
				Type:     spec.TYPE_CATEGORY,
				Title:    v.Name,
				Items:    walkCollection(v.Items, id),
			}
			cs = append(cs, specItem)
		}
	}
	return cs
}

func convertContent(item Item) []*spec.CollectionNode {
	req := spec.NewCollectionHttpRequest()
	for k, v := range item.Request.Url.Path {
		if !strings.HasPrefix(v, ":") {
			continue
		}
		for _, x := range item.Request.Url.Variables {
			if x.Key == v[1:] {
				item.Request.Url.Path[k] = "{" + x.Key + "}"
				req.Attrs.Parameters.Path = append(req.Attrs.Parameters.Path, x.toParameter())
				break
			}
		}
	}

	url := spec.NewCollectionHttpUrl("/"+strings.Join(item.Request.Url.Path, "/"), item.Request.Method)
	nodes := spec.CollectionNodes{
		url.ToCollectionNode(),
	}

	for _, v := range item.Request.Url.Queries {
		if v.Disabled {
			continue
		}
		req.Attrs.Parameters.Query = append(req.Attrs.Parameters.Query, v.toParameter())
	}
	for _, v := range item.Request.Headers {
		req.Attrs.Parameters.Header = append(req.Attrs.Parameters.Header, v.toParameter())
	}

	if body := encodeRequestBody(item.Request.Body); body != nil {
		req.Attrs.Content = body
	} else {
		req.Attrs.Content = spec.HTTPBody{
			"application/json": {Schema: jsonschema.NewSchema("object")},
		}
	}
	nodes = append(nodes, req.ToCollectionNode())

	res := spec.NewCollectionHttpResponse()
	res.Attrs = encodeResponseBody(item.Response)
	nodes = append(nodes, res.ToCollectionNode())
	return nodes
}

var contenttypemapp = map[string]string{
	"json":      "application/json",
	"urlencode": "application/x-www-form-urlencoded",
	"formdata":  "multipart/form-data",
	"plain":     "text/plain",
}

func encodeRequestBody(body *Body) spec.HTTPBody {
	if body == nil || body.Disabled {
		return nil
	}
	switch body.Mode {
	case "raw":
		if body.Options.Raw.Language == "json" {
			b := jsonToSchema(body.Raw)
			return map[string]*spec.Body{
				contenttypemapp["json"]: {
					Schema: b,
				},
			}
		}
	case "formdata", "urlencode":
		b := jsonschema.NewSchema("object")
		b.Properties = make(map[string]*jsonschema.Schema)
		for _, v := range body.Formdata {
			if v.Disabled {
				continue
			}
			b.Properties[v.Key] = v.toJSONSchema()
		}
		return map[string]*spec.Body{
			contenttypemapp[body.Mode]: {
				Schema: b,
			},
		}
	case "file":
	case "graphql":
	default:
	}

	return nil
}

func encodeResponseBody(res []Response) *spec.HttpResponseAttrs {
	response := &spec.HttpResponseAttrs{
		List: make(spec.Responses, 0),
	}
	for _, v := range res {
		r := spec.Response{Code: v.Code}
		r.Description = v.Name
		switch v.PostmanePreviewLanguage {
		case "json":

			// fmt.Println("json.........")
			b := jsonToSchema(v.Body)
			b.Examples = v.Body
			r.Content = map[string]*spec.Body{
				contenttypemapp["json"]: {
					Schema: b,
				},
			}

		case "plain":
			b := jsonschema.NewSchema(jsonschema.T_STR)
			b.Examples = v.Body
			r.Content = map[string]*spec.Body{
				contenttypemapp["plain"]: {
					Schema: b,
				},
			}
		default:
			r.Content = map[string]*spec.Body{
				contenttypemapp["plain"]: {
					Schema: jsonschema.NewSchema(jsonschema.T_OBJ),
				},
			}
		}

		for _, v := range v.Header {
			if v.Disabled {
				continue
			}
			b := jsonschema.NewSchema(jsonschema.T_STR)
			b.Examples = v.Value
			r.Header = append(r.Header, &spec.Parameter{
				Name:        v.Key,
				Description: v.Description,
				Schema:      b,
			})
		}

		response.List = append(response.List, &r)
	}

	if len(response.List) == 0 {
		defaultres := spec.Response{Code: 200}
		defaultres.Name = "success"
		defaultres.Content = spec.HTTPBody{
			"application/json": {Schema: jsonschema.NewSchema(jsonschema.T_OBJ)},
		}
		response.List = append(response.List, &defaultres)
	}
	return response
}
