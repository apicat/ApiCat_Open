package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec2"
	"github.com/apicat/apicat/v2/backend/module/spec2/jsonschema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

type openapiParser struct {
	modelMapping      map[string]int64
	parametersMapping map[string]*spec2.Parameter
}
type openapiGenerator struct {
	modelMapping map[int64]string
}

type openapiSpec struct {
	Openapi    string                                `json:"openapi"`
	Info       spec2.Info                            `json:"info"`
	Servers    []spec2.Server                        `json:"servers,omitempty"`
	Components map[string]any                        `json:"components,omitempty"`
	Paths      map[string]map[string]openapiPathItem `json:"paths"`
	Tags       []tagObject                           `json:"tags,omitempty"`
}

type openapiPathItem struct {
	Summary     string              `json:"summary"`
	Description string              `json:"description,omitempty"`
	OperationId string              `json:"operationId"`
	Tags        []string            `json:"tags,omitempty"`
	Parameters  []openAPIParamter   `json:"parameters,omitempty"`
	RequestBody *openapiRequestbody `json:"requestBody,omitempty"`
	Responses   map[string]any      `json:"responses,omitempty"`
}

type openapiRequestbody struct {
	Content map[string]*jsonschema.Schema `json:"content,omitempty"`
}

func (o *openapiParser) parseInfo(info *base.Info) spec2.Info {
	return spec2.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (o *openapiParser) parseServers(servs []*v3.Server) []spec2.Server {
	srvs := make([]spec2.Server, len(servs))
	for k, v := range servs {
		srvs[k] = spec2.Server{
			URL:         v.URL,
			Description: v.Description,
		}
	}
	return srvs
}

func (o *openapiParser) parseContent(mts map[string]*v3.MediaType) (spec2.HTTPBody, error) {
	if mts == nil {
		return nil, errors.New("no content")
	}

	content := make(spec2.HTTPBody)
	for contentType, mediaType := range mts {
		body := &spec2.Body{}

		js, err := jsonSchemaConverter(mediaType.Schema)
		if err != nil {
			return nil, err
		}
		js.Examples = mediaType.Example

		if len(mediaType.Examples) > 0 {
			body.Examples = make([]spec2.Example, 0)
			for _, v := range mediaType.Examples {
				if example, err := json.Marshal(v); err == nil {
					body.Examples = append(body.Examples, spec2.Example{
						Summary: v.Summary,
						Value:   string(example),
					})
				}
			}
		}
		body.Schema = js
		content[contentType] = body
	}
	return content, nil
}

func (o *openapiParser) parseDefinetions(comp *v3.Components) (*spec2.Definitions, error) {
	if comp == nil {
		return &spec2.Definitions{
			Schemas:   make(spec2.DefinitionModels, 0),
			Responses: make(spec2.DefinitionResponses, 0),
		}, nil
	}

	o.modelMapping = map[string]int64{}
	models := make(spec2.DefinitionModels, 0)

	for k, v := range comp.Schemas {
		js, err := jsonSchemaConverter(v)
		if err != nil {
			return nil, err
		}
		models = append(models, &spec2.DefinitionModel{
			ID:          stringToUnid(k),
			Name:        k,
			Description: js.Description,
			Schema:      js,
		})
	}

	responses := make(spec2.DefinitionResponses, 0)
	for k, v := range comp.Responses {
		id := stringToUnid(k)
		def := &spec2.DefinitionResponse{
			BasicResponse: spec2.BasicResponse{
				Header:      make(spec2.ParameterList, 0),
				Name:        k,
				ID:          id,
				Description: v.Description,
			},
		}

		if v.Headers != nil {
			for k, v := range v.Headers {
				js, err := jsonSchemaConverter(v.Schema)
				if err != nil {
					return nil, err
				}
				js.Description = v.Description
				def.Header = append(def.Header, &spec2.Parameter{
					Name:   k,
					Schema: js,
				})
			}
		}
		if v.Content != nil {
			content, err := o.parseContent(v.Content)
			if err != nil {
				return nil, err
			}
			def.Content = content
		}
		responses = append(responses, def)
	}
	return &spec2.Definitions{
		Schemas:   models,
		Responses: responses,
	}, nil
}

func (o *openapiParser) parseGlobalParameters(com *v3.Components) *spec2.GlobalParameters {
	res := spec2.NewGlobalParameters()
	if com == nil {
		return res
	}

	inp := com.Extensions
	if inp == nil {
		return res
	}

	global, ok := inp["x-apicat-global-parameters"]
	if !ok {
		return res
	}

	for k, v := range global.(map[string]any) {
		nb, err := json.Marshal(v)
		if err != nil {
			continue
		}

		s := &spec2.Parameter{}
		json.Unmarshal(nb, s)
		in := strings.Index(k, "-")
		if in == -1 {
			continue
		}
		res.Add(k[:in], s)
	}
	return res
}

func (o *openapiParser) parseParameters(inp []*v3.Parameter) (*spec2.HTTPParameters, error) {
	rawparamter := spec2.NewHTTPParameters()
	for _, v := range inp {
		if g := v.GoLow(); g.IsReference() {
			// if this parameter is a global parameter
			if isGlobalParameter(g.Reference.Reference) {
				continue
			}

			name := fmt.Sprintf("%s-%s", v.In, v.Name)
			sc, ok := o.parametersMapping[name]
			if ok {
				rawparamter.Add(v.In, sc)
				continue
				// r := fmt.Sprintf("#/definitions/parameters/%d", id)
				// rawparamter.Add(v.In, &spec.Schema{
				// 	Reference: &r,
				// })
				// continue
			}
		}

		sp := &spec2.Parameter{
			Name:     v.Name,
			Required: v.Required,
		}

		sp.Schema = &jsonschema.Schema{}
		if v.Schema != nil {
			js, err := jsonSchemaConverter(v.Schema)
			if err != nil {
				return nil, err
			}
			sp.Schema = js
		}
		sp.Schema.Description = v.Description
		sp.Schema.Examples = v.Example
		sp.Schema.Deprecated = v.Deprecated
		rawparamter.Add(v.In, sp)
	}
	return rawparamter, nil
}

func (o *openapiParser) parseResponses(responses map[string]*v3.Response) (*spec2.CollectionHttpResponse, error) {
	outresponses := spec2.NewCollectionHttpResponse()

	for code, res := range responses {
		c, _ := strconv.Atoi(code)
		resp := spec2.Response{
			Code: c,
		}

		s := res.GoLow().Reference.Reference
		if s != "" {
			refs := fmt.Sprintf("#/definitions/responses/%d", stringToUnid(s[strings.LastIndex(s, "/")+1:]))
			resp.Reference = refs
			outresponses.Attrs.List = append(outresponses.Attrs.List, &resp)
			continue
		}

		if _, ok := res.Extensions["x-apicat-response-name"]; ok {
			resp.Name = res.Extensions["x-apicat-response-name"].(string)
		}

		resp.Description = res.Description
		resp.Header = make(spec2.ParameterList, 0)
		if res.Headers != nil {
			for k, v := range res.Headers {
				js, err := jsonSchemaConverter(v.Schema)
				if err != nil {
					return nil, err
				}

				js.Description = v.Description
				resp.Header = append(resp.Header, &spec2.Parameter{
					Name:   k,
					Schema: js,
				})
			}
		}

		content, err := o.parseContent(res.Content)
		if err != nil {
			return nil, err
		}
		resp.Content = content
		outresponses.Attrs.List = append(outresponses.Attrs.List, &resp)
	}
	return outresponses, nil
}

func (o *openapiParser) parseCollections(paths *v3.Paths) spec2.Collections {
	collections := make(spec2.Collections, 0)
	if paths == nil {
		return collections
	}

	var err error
	for path, p := range paths.PathItems {
		op := p.GetOperations()
		for method, info := range op {
			content := spec2.CollectionNodes{
				spec2.NewCollectionHttpUrl(path, method).ToCollectionNode(),
			}

			// parse markdown to doc
			// doctree := markdown.ToDocment([]byte(info.Description))
			// for _, v := range doctree.Items {
			// 	content = append(content, spec.MuseCreateNodeProxy(v))
			// }

			// request
			req := spec2.NewCollectionHttpRequest()
			if req.Attrs.Parameters, err = o.parseParameters(info.Parameters); err != nil {
				continue
			}
			if info.RequestBody != nil {
				if req.Attrs.Content, err = o.parseContent(info.RequestBody.Content); err != nil {
					continue
				}
			}
			content = append(content, req.ToCollectionNode())

			// response
			res := spec2.NewCollectionHttpResponse()
			if info.Responses != nil {
				if res, err = o.parseResponses(info.Responses.Codes); err != nil {
					continue
				}
			}
			content = append(content, res.ToCollectionNode())

			title := info.Summary
			if title == "" {
				title = path
			}

			collections = append(collections, &spec2.Collection{
				Type:    spec2.TYPE_HTTP,
				Title:   title,
				Tags:    info.Tags,
				Content: content,
			})
		}
	}
	return collections
}

func (o *openapiGenerator) generateBase(in *spec2.Spec, version string) *openapiSpec {
	return &openapiSpec{
		Openapi: version,
		Info: spec2.Info{
			Title:       in.Info.Title,
			Description: in.Info.Description,
			Version:     in.Info.Version,
		},
		Servers:    in.Servers,
		Components: o.generateComponents(version, in),
	}
}

func (o *openapiGenerator) convertJsonSchema(version string, in *jsonschema.Schema) *jsonschema.Schema {
	if in == nil {
		return nil
	}
	p := convertJsonSchemaRef(in, version, o.modelMapping)
	// if p.Reference == nil && strings.HasPrefix(ver, "3.0") {
	// 	if p.Items != nil {
	// 		if !p.Items.IsBool() {
	// 			p.Items.SetValue(&jsonschema.Schema{})
	// 		}
	// 	}
	// 	if p.Properties != nil {
	// 		for _, v := range p.Properties {
	// 			o.convertJSONSchema(ver, v)
	// 		}
	// 	}
	// 	if p.AdditionalProperties != nil {
	// 		if !p.AdditionalProperties.IsBool() {
	// 			o.convertJSONSchema(ver, p.AdditionalProperties.Value())
	// 		}
	// 	}
	// 	p.Type.SetValue(p.Type.Value()[0])
	// }
	if p.Type != nil {
		t := p.Type.List()
		if len(t) > 0 && t[0] == "file" {
			// jsonschema 没有file
			p.Type.Set(jsonschema.T_ARR)
			p.Items = &jsonschema.ValueOrBoolean[*jsonschema.Schema]{}
			p.Items.SetValue(&jsonschema.Schema{})
		}
	}
	return p
}

func (o *openapiGenerator) generateResponseWithoutRef(resp *spec2.BasicResponse, version string) map[string]any {
	result := map[string]any{}
	if resp.Content != nil {
		c := make(map[string]*jsonschema.Schema)
		for contentType, body := range resp.Content {
			jsonSchema := o.convertJsonSchema(version, body.Schema)
			c[contentType] = jsonSchema
		}
		result["content"] = c
	}

	if len(resp.Header) > 0 {
		headers := make(map[string]any)
		for _, h := range resp.Header {
			headers[h.Name] = map[string]any{
				"description": h.Description,
				"schema":      o.convertJsonSchema(version, h.Schema),
			}
		}
		result["headers"] = headers
	}

	result["description"] = resp.Description
	result["x-apicat-response-name"] = resp.Name
	return result
}

func (o *openapiGenerator) generateReqParams(collectionReq spec2.CollectionHttpRequest, globalsParmaters *spec2.GlobalParameters, version string) []openAPIParamter {
	// var out []openAPIParamter
	out := globalToLocalParameters(globalsParmaters, false, collectionReq.Attrs.GlobalExcepts.ToMap())

	for in, params := range collectionReq.Attrs.Parameters.ToMap() {
		for _, param := range params {
			p := *param
			// if p.Reference != nil {
			// 	if defp := spe.Definitions.Parameters.LookupID(toInt64(getRefName(*p.Reference))); defp != nil {
			// 		p = *defp
			// 	}
			// }
			item := openAPIParamter{
				Name:        p.Name,
				Required:    p.Required,
				Description: p.Description,
				// Example:     p.Example,
				Schema: o.convertJsonSchema(version, p.Schema),
				In:     in,
			}
			out = append(out, item)
		}
	}
	return out
}

func (o *openapiGenerator) generateResponse(resp *spec2.Response, definitionsResps spec2.DefinitionResponses, version string) map[string]any {
	if resp.Reference != "" {
		if strings.HasPrefix(resp.Reference, "#/definitions/responses/") {
			if x := definitionsResps.FindByID(
				toInt64(getRefName(resp.Reference)),
			); x != nil {
				name_id := fmt.Sprintf("%s-%d", x.Name, x.ID)
				return map[string]any{
					"$ref": "#/components/responses/" + name_id,
				}
			}
		}
	}
	return o.generateResponseWithoutRef(&resp.BasicResponse, version)
}

func (o *openapiGenerator) generateComponents(version string, in *spec2.Spec) map[string]any {
	schemas := make(map[string]jsonschema.Schema)
	respons := make(map[string]any)
	o.modelMapping = map[int64]string{}

	definitionModels := spec2.DefinitionModels{}
	for _, v := range in.Definitions.Schemas {
		if v.Type == string(spec2.TYPE_CATEGORY) {
			items := v.ItemsTreeToList()
			for _, item := range items {
				o.modelMapping[item.ID] = item.Name
			}
			definitionModels = append(definitionModels, items...)
		} else {
			o.modelMapping[v.ID] = v.Name
			definitionModels = append(definitionModels, v)
		}
	}
	for _, v := range definitionModels {
		name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(v.Name, " ", ""), v.ID)
		schemas[name_id] = *o.convertJsonSchema(version, v.Schema)
	}

	for _, v := range in.Definitions.Responses {
		if v.Type == string(spec2.TYPE_CATEGORY) {
			resps := v.ItemsTreeToList()
			for _, resp := range resps {
				name_id := fmt.Sprintf("%s-%d", resp.Name, resp.ID)
				respons[name_id] = o.generateResponseWithoutRef(&resp.BasicResponse, version)
			}
		} else {
			name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(v.Name, " ", ""), v.ID)
			respons[name_id] = o.generateResponseWithoutRef(&v.BasicResponse, version)
		}
	}

	globalParam := in.Globals.Parameters
	m := globalParam.ToMap()
	globals := make(map[string]openAPIParamter)
	for in, ps := range m {
		for _, p := range ps {
			globals[fmt.Sprintf("%s-%s", in, strings.ReplaceAll(p.Name, " ", ""))] = toParameter(p, in, version)
		}
	}

	return map[string]any{
		"schemas":                    schemas,
		"responses":                  respons,
		"x-apicat-global-parameters": globals,
	}
}

func (o *openapiGenerator) generatePaths(version string, in *spec2.Spec) (map[string]map[string]openapiPathItem, []tagObject) {
	var (
		out  = make(map[string]map[string]openapiPathItem)
		tags = make(map[string]struct{})
	)

	for path, ops := range deepGetHttpCollection(&in.Collections) {
		if path == "" {
			continue
		}

		for method, op := range ops {
			item := openapiPathItem{
				Summary:     op.Title,
				Description: op.Description,
				OperationId: op.OperatorID,
				Tags:        op.Tags,
				Parameters:  o.generateReqParams(op.Req, in.Globals.Parameters, version),
				Responses:   make(map[string]any),
			}

			for _, v := range op.Tags {
				tags[v] = struct{}{}
			}

			for contentType, body := range op.Req.Attrs.Content {
				if contentType == "none" {
					continue
				}

				sp := o.convertJsonSchema(version, body.Schema)
				sp.Examples = body.Examples
				if item.RequestBody == nil {
					item.RequestBody = &openapiRequestbody{
						Content: make(map[string]*jsonschema.Schema),
					}
				}
				item.RequestBody.Content[contentType] = sp
			}

			for _, v := range op.Res.Attrs.List {
				res := o.generateResponse(v, in.Definitions.Responses, version)
				item.Responses[strconv.Itoa(v.Code)] = res
			}

			if len(op.Res.Attrs.List) == 0 {
				item.Responses["200"] = map[string]any{
					"description": "success",
				}
			}

			if _, ok := out[path]; !ok {
				out[path] = make(map[string]openapiPathItem)
			}

			out[path][method] = item
		}
	}

	return out, func() (list []tagObject) {
		for k := range tags {
			list = append(list, tagObject{Name: k})
		}
		return
	}()
}
