package openapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/apicat/apicat/v2/backend/module/spec"
	"github.com/apicat/apicat/v2/backend/module/spec2"
	"github.com/apicat/apicat/v2/backend/module/spec2/jsonschema"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v2 "github.com/pb33f/libopenapi/datamodel/high/v2"
)

type swaggerParser struct {
	modelMapping           map[string]int64
	parametersMapping      map[string]*spec2.Parameter
	globalParamtersMapping map[string]struct{}
}

type swaggerGenerator struct {
	modelNames map[int64]string
}

type swaggerSpec struct {
	Swagger          string                                `json:"swagger"`
	Info             *spec2.Info                           `json:"info"`
	Tags             []tagObject                           `json:"tags,omitempty"`
	Host             string                                `json:"host,omitempty"`
	BasePath         string                                `json:"basePath"`
	Schemas          []string                              `json:"schemes,omitempty"`
	Definitions      map[string]jsonschema.Schema          `json:"definitions"`
	Parameters       map[string]openAPIParamter            `json:"parameters,omitempty"`
	Responses        map[string]any                        `json:"responses,omitempty"`
	Paths            map[string]map[string]swaggerPathItem `json:"paths"`
	GlobalParameters map[string]openAPIParamter            `json:"x-apicat-global-parameters,omitempty"`
}

type swaggerPathItem struct {
	Summary     string            `json:"summary"`
	Tags        []string          `json:"tags,omitempty"`
	Description string            `json:"description,omitempty"`
	OperationId string            `json:"operationId"`
	Consumes    []string          `json:"consumes,omitempty"`
	Produces    []string          `json:"produces,omitempty"`
	Parameters  []openAPIParamter `json:"parameters,omitempty"`
	Responses   map[string]any    `json:"responses,omitempty"`
}

func (s *swaggerParser) parseInfo(info *base.Info) spec2.Info {
	return spec2.Info{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
	}
}

func (s *swaggerParser) parseServers(in *v2.Swagger) []spec2.Server {
	servers := make([]spec2.Server, len(in.Schemes))
	if in.BasePath == "/" {
		in.BasePath = ""
	}
	for k, v := range in.Schemes {
		servers[k] = spec2.Server{
			URL:         fmt.Sprintf("%s://%s%s", v, in.Host, in.BasePath),
			Description: v,
		}
	}
	return servers
}

func (s *swaggerParser) parseDefinitionModels(defs *v2.Definitions) (spec2.DefinitionModels, error) {
	s.modelMapping = make(map[string]int64)
	models := make(spec2.DefinitionModels, 0)
	if defs == nil {
		return models, nil
	}

	for k, v := range defs.Definitions {
		js, err := jsonSchemaConverter(v)
		if err != nil {
			return nil, err
		}

		id := stringToUnid(k)
		s.modelMapping[k] = id
		models = append(models, &spec2.DefinitionModel{
			ID:          id,
			Name:        k,
			Description: k,
			Schema:      js,
		})
	}

	return models, nil
}

func (s *swaggerParser) parseJsonSchema(b *base.SchemaProxy) (*jsonschema.Schema, error) {
	js, err := jsonSchemaConverter(b)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func (s *swaggerParser) parseDefinitionResponses(in *v2.Swagger) (spec2.DefinitionResponses, error) {
	list := make(spec2.DefinitionResponses, 0)
	if in.Responses == nil {
		return list, nil
	}

	for key, res := range in.Responses.Definitions {
		headers := make([]*spec2.Parameter, 0)
		content := make(spec2.HTTPBody)
		if res.Headers != nil {
			for k, v := range res.Headers {
				headers = append(headers, &spec2.Parameter{
					Name: k,
					Schema: &jsonschema.Schema{
						Type:        jsonschema.NewSchemaType(v.Type),
						Format:      v.Format,
						Description: v.Description,
						Examples:    v.Default,
					},
				})
			}
		}
		if res.Schema != nil {
			js, err := s.parseJsonSchema(res.Schema)
			if err != nil {
				return list, err
			}

			body := &spec2.Body{Schema: js}
			if len(in.Produces) == 0 {
				content["application/json"] = body
				if res.Examples != nil {
					body.Examples = make([]spec2.Example, 0)
					for k, v := range res.Examples.Values {
						if example, err := json.Marshal(v); err == nil {
							body.Examples = append(body.Examples, spec2.Example{
								Summary: k,
								Value:   string(example),
							})
						}
					}
				}
			} else {
				for _, v := range in.Produces {
					content[v] = body
					if res.Examples != nil {
						emp, ok := res.Examples.Values[v]
						if ok {
							if example, err := json.Marshal(emp); err == nil {
								body.Examples = append(body.Examples, spec2.Example{
									Summary: v,
									Value:   string(example),
								})
							}
						}
					}
				}
			}
		}
		list = append(list, &spec2.DefinitionResponse{
			BasicResponse: spec2.BasicResponse{
				ID:          stringToUnid(key),
				Name:        key,
				Header:      headers,
				Content:     content,
				Description: res.Description,
			},
		})
	}
	return list, nil
}

func (s *swaggerParser) parseGlobalParameters(inp map[string]any) *spec2.GlobalParameters {
	params := spec2.NewGlobalParameters()
	if inp == nil {
		return params
	}
	global, ok := inp["x-apicat-global-parameters"]
	if !ok {
		return params
	}

	s.globalParamtersMapping = make(map[string]struct{})

	for k, v := range global.(map[string]any) {
		nb, err := json.Marshal(v)
		if err != nil {
			continue
		}

		p := &spec2.Parameter{}
		json.Unmarshal(nb, p)
		in := strings.Index(k, "-")
		if in == -1 {
			continue
		}
		params.Add(k[:in], p)
		s.globalParamtersMapping[p.Name] = struct{}{}
	}
	return params
}

func (s *swaggerParser) parseRequest(in *v2.Swagger, info *v2.Operation) (*spec2.CollectionHttpRequest, error) {
	request := spec2.NewCollectionHttpRequest()

	var err error
	body := &spec2.Body{}
	// 有效载荷application/x-www-form-urlencoded和multipart/form-data请求是通过使用form参数来描述，而不是body参数。
	formData := &jsonschema.Schema{
		Type:       jsonschema.NewSchemaType(jsonschema.T_OBJ),
		Properties: make(map[string]*jsonschema.Schema),
	}

	for _, v := range info.Parameters {
		// 这里引用 #/parameters 暂时无法获取
		// 直接展开
		if _, ok := s.parametersMapping[v.Name]; ok {
			continue
		}
		if _, ok := s.globalParamtersMapping[v.Name]; ok {
			continue
		}

		required := v.Required != nil && *v.Required
		switch v.In {
		case "query", "header", "path", "cookie":
			request.Attrs.Parameters.Add(v.In,
				&spec2.Parameter{
					Name:        v.Name,
					Description: v.Description,
					Required:    required,
					Schema: &jsonschema.Schema{
						Type:   jsonschema.NewSchemaType(v.Type),
						Format: v.Format,
					},
				},
			)
		case "formData":
			formData.Properties[v.Name] = &jsonschema.Schema{
				Type:        jsonschema.NewSchemaType(v.Type),
				Description: v.Description,
				Format:      v.Format,
				Default:     v.Default,
			}
			if required {
				formData.Required = append(formData.Required, v.Name)
			}
		case "body":
			body.Schema, err = s.parseJsonSchema(v.Schema)
			if err != nil {
				return nil, err
			}
		}
	}

	consumes := info.Consumes
	if len(info.Consumes) == 0 {
		// 从global获取
		consumes = in.Consumes
	}
	// 有些文件没有consunmer 给个默认 否则body不知道什么是mine
	// if len(consumes) == 0 && body != nil {
	// 	consumes = []string{defaultSwaggerConsumerProduce}
	// }

	for _, v := range consumes {
		if strings.Contains(v, "form") {
			request.Attrs.Content[v] = &spec2.Body{Schema: formData}
		} else {
			if body.Schema != nil {
				request.Attrs.Content[v] = body
			}
		}
	}
	return request, nil
}

func (s *swaggerParser) parseResponse(info *v2.Operation) (*spec2.CollectionHttpResponse, error) {
	if info.Responses == nil {
		return nil, nil
	}
	outresponses := spec2.NewCollectionHttpResponse()
	// if info.Responses.Default != nil {
	// 	// 我们没有default
	// 	// todo
	// }
	for code, res := range info.Responses.Codes {
		// res github.com/pb33f/libopenapi 不支持response ref 所以无法获取
		// 这里的common无法转换
		c, err := strconv.Atoi(code)
		if err != nil {
			continue
		}
		resp := spec2.Response{
			Code: c,
		}
		if _, ok := res.Extensions["x-apicat-response-name"]; ok {
			resp.Name = res.Extensions["x-apicat-response-name"].(string)
		}
		resp.Description = res.Description
		resp.Content = make(spec2.HTTPBody)
		resp.Header = make(spec2.ParameterList, 0)

		// libopenapi not support response ref, in swagger 2.0
		// it's like dereference
		if res.GoLow().Schema.GetReference() != "" {
			ref := res.GoLow().Schema.GetReference()
			refs := fmt.Sprintf("#/definitions/responses/%d", stringToUnid(ref[strings.LastIndex(ref, "/")+1:]))
			resp.Reference = refs
			outresponses.Attrs.List = append(outresponses.Attrs.List, &resp)
			continue
		}
		if res.Headers != nil {
			for k, v := range res.Headers {
				resp.Header = append(resp.Header, &spec2.Parameter{
					Name: k,
					Schema: &jsonschema.Schema{
						Type:        jsonschema.NewSchemaType(v.Type),
						Format:      v.Format,
						Description: v.Description,
						Examples:    v.Default,
					},
				})
			}
		}
		if res.Schema != nil {
			js, err := s.parseJsonSchema(res.Schema)
			if err != nil {
				return nil, err
			}
			for _, v := range info.Produces {
				body := &spec2.Body{Schema: js}
				if res.Examples != nil {
					mp, ok := res.Examples.Values[v]
					if ok {
						if example, err := json.Marshal(mp); err == nil {
							body.Examples = append(body.Examples, spec2.Example{
								Summary: v,
								Value:   string(example),
							})
						}
					}
				}
				resp.Content[v] = body
			}
		}
		outresponses.Attrs.List = append(outresponses.Attrs.List, &resp)
	}
	if len(outresponses.Attrs.List) == 0 {
		outresponses.Attrs.List = append(outresponses.Attrs.List, &spec2.Response{
			Code:          200,
			BasicResponse: spec2.BasicResponse{Description: "success"},
		})
	}
	return outresponses, nil
}

func (s *swaggerParser) parseCollections(in *v2.Swagger, paths *v2.Paths) spec2.Collections {
	collections := make(spec2.Collections, 0)
	for path, p := range paths.PathItems {
		op := p.GetOperations()
		for method, info := range op {
			content := spec2.CollectionNodes{
				spec2.NewCollectionHttpUrl(path, method).ToCollectionNode(),
			}

			// parse markdown to doc
			// doctree := markdown.ToDocment([]byte(info.Description))
			// for _, v := range doctree.Items {
			// 	content = append(content, spec2.NewCollectionDoc(v).ToCollectionNode())
			// }

			// request
			if req, err := s.parseRequest(in, info); err != nil {
				continue
			} else {
				content = append(content, req.ToCollectionNode())
			}

			// response
			if res, err := s.parseResponse(info); err != nil {
				continue
			} else {
				content = append(content, res.ToCollectionNode())
			}

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

func (s *swaggerGenerator) convertJsonSchema(v *jsonschema.Schema) *jsonschema.Schema {
	if v == nil {
		return v
	}
	return convertJsonSchemaRef(v, "2.0", s.modelNames)
}

func (s *swaggerGenerator) generateBase(in *spec2.Spec) *swaggerSpec {
	s.modelNames = map[int64]string{}
	out := &swaggerSpec{
		Swagger: "2.0",
		Info: &spec2.Info{
			Title:       in.Info.Title,
			Description: in.Info.Description,
			Version:     in.Info.Version,
		},
		Definitions: make(map[string]jsonschema.Schema),
	}

	for _, v := range in.Servers {
		u, err := url.Parse(v.URL)
		if err != nil {
			continue
		}

		if out.Host == "" {
			out.Host = u.Host
			out.BasePath = u.Path
		}
		out.Schemas = append(out.Schemas, u.Scheme)
		// just need fist one
		break
	}

	definitionModels := spec2.DefinitionModels{}
	for _, v := range in.Definitions.Schemas {
		if v.Type == string(spec2.TYPE_CATEGORY) {
			items := v.ItemsTreeToList()
			for _, item := range items {
				s.modelNames[item.ID] = item.Name
			}
			definitionModels = append(definitionModels, items...)
		} else {
			s.modelNames[v.ID] = v.Name
			definitionModels = append(definitionModels, v)
		}
	}

	for _, v := range definitionModels {
		name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(v.Name, " ", ""), v.ID)
		out.Definitions[name_id] = *s.convertJsonSchema(v.Schema)
	}

	globalParams := in.Globals.Parameters.ToMap()
	out.GlobalParameters = make(map[string]openAPIParamter)
	for in, paramList := range globalParams {
		for _, p := range paramList {
			name_id := fmt.Sprintf("%s-%d", p.Name, p.ID)
			out.GlobalParameters[name_id] = toParameter(p, in, "2.0")
		}
	}

	if out.BasePath == "" {
		out.BasePath = "/"
	}
	if len(in.Definitions.Responses) > 0 {
		out.Responses = make(map[string]any)
		for _, v := range in.Definitions.Responses {
			if v.Type == string(spec2.TYPE_CATEGORY) {
				items := v.ItemsTreeToList()
				for _, item := range items {
					name_id := fmt.Sprintf("%s-%d", item.Name, item.ID)
					out.Responses[name_id] = s.generateResponseWithoutRef(in, &item.BasicResponse)
				}
			} else {
				name_id := fmt.Sprintf("%s-%d", strings.ReplaceAll(v.Name, " ", ""), v.ID)
				out.Responses[name_id] = s.generateResponseWithoutRef(in, &v.BasicResponse)
			}
		}
	}

	return out
}

func (s *swaggerGenerator) generateReqParams(collectionReq spec2.CollectionHttpRequest, spe *spec2.Spec) []openAPIParamter {
	// 添加启用的全局参数
	out := globalToLocalParameters(spe.Globals.Parameters, true, collectionReq.Attrs.GlobalExcepts.ToMap())

	for in, params := range collectionReq.Attrs.Parameters.ToMap() {
		switch in {
		case "header", "query", "path", "cookie":
			for _, v := range params {
				// if v.Reference != nil {
				// 	// 解开公共参数
				// 	if id := toInt64(getRefName(*v.Reference)); id != 0 {
				// 		v = spe.Definitions.Parameters.LookupID(id)
				// 	}
				// }
				newv := *v
				newv.Schema = s.convertJsonSchema(v.Schema)
				out = append(out, toParameter(&newv, in, "2.0"))
			}
		}
	}

	if collectionReq.Attrs.Content == nil {
		return out
	}

	var hasBody bool
	for contentType, body := range collectionReq.Attrs.Content {
		// contentType incloud form use parameters in
		if strings.Contains(contentType, "form") {
			if body.Schema == nil {
				continue
			}

			if num := len(body.Schema.Type.List()); num == 0 {
				continue
			}

			typ := body.Schema.Type.First()
			if typ != jsonschema.T_OBJ || body.Schema.Properties == nil {
				continue
			}

			for k, v := range body.Schema.Properties {
				content := openAPIParamter{
					Name:        k,
					In:          "formData",
					Type:        v.Type.First(),
					Description: v.Description,
					Schema:      s.convertJsonSchema(v),
					Required: func() bool {
						for _, r := range v.Required {
							if r == k {
								return true
							}
						}
						return false
					}(),
				}
				if v != nil {
					t := v.Type.List()
					if len(t) > 0 && t[0] == "file" {
						content.Type = t[0]
					}
				}
				out = append(out, content)
			}
		} else {
			if hasBody {
				continue
			}

			out = append(out, openAPIParamter{
				Name:     "body",
				Schema:   s.convertJsonSchema(body.Schema),
				In:       "body",
				Required: true,
			})
			hasBody = true
		}
	}
	return out
}

func (s *swaggerGenerator) generateResponseWithoutRef(in *spec2.Spec, resp *spec2.BasicResponse) map[string]any {
	response := map[string]any{
		"x-apicat-response-name": resp.Name,
		"description":            resp.Description,
	}

	if len(resp.Header) > 0 {
		header := make(map[string]any)
		for _, v := range resp.Header {
			if v.Schema.Description == "" {
				v.Schema.Description = v.Description
			}
			v.Schema.Default = v.Schema.Examples
			v.Schema.Examples = nil
			header[v.Name] = v.Schema
		}
		response["headers"] = header
	}

	if resp.Content != nil {
		for k, v := range resp.Content {
			response["schema"] = s.convertJsonSchema(v.Schema)
			if v.Examples != nil {
				for _, v := range v.Examples {
					response["examples"] = map[string]any{
						k: v,
					}
					break
				}
			}
			break
		}
	}
	return response
}

func (s *swaggerGenerator) generateResponse(in *spec2.Spec, resp *spec2.Response) map[string]any {
	if resp.Reference != "" {
		if strings.HasPrefix(resp.Reference, "#/definitions/responses/") {
			x := in.Definitions.Responses.FindByID(
				toInt64(getRefName(resp.Reference)),
			)
			if x != nil {
				name := fmt.Sprintf("%s-%d", x.Name, x.ID)
				return map[string]any{
					"$ref": "#/responses/" + name,
				}
			}
		}
		return nil
	}
	return s.generateResponseWithoutRef(in, &resp.BasicResponse)
}

func (s *swaggerGenerator) generatePathResponse(in *spec2.Spec, resp spec2.CollectionHttpResponse) (map[string]any, []string) {
	product := map[string]struct{}{}
	result := make(map[string]any)

	for _, r := range resp.Attrs.List {
		result[strconv.Itoa(r.Code)] = s.generateResponse(in, r)
		for k := range r.Content {
			if _, ok := product[k]; !ok {
				product[k] = struct{}{}
				continue
			}
		}
	}
	if len(result) == 0 {
		result["default"] = map[string]string{
			"description": "success",
		}
	}
	return result, func() (ret []string) {
		if len(product) == 0 {
			return []string{"application/json"}
		}
		for k := range product {
			ret = append(ret, k)
		}
		return
	}()
}

func (s *swaggerGenerator) generatePaths(in *spec2.Spec) (map[string]map[string]swaggerPathItem, []tagObject) {
	out := make(map[string]map[string]swaggerPathItem)
	tags := make(map[string]struct{})

	for path, methods := range deepGetHttpCollection(&in.Collections) {
		if path == "" {
			continue
		}
		for method, op := range methods {
			reslist, product := s.generatePathResponse(in, op.Res)
			if len(reslist) == 0 {
				reslist["default"] = &spec.Schema{Description: "success"}
			}
			item := swaggerPathItem{
				Summary:     op.Title,
				Description: op.Description,
				OperationId: op.OperatorID,
				Parameters:  s.generateReqParams(op.Req, in),
				Produces:    product,
				Responses:   reslist,
				Tags:        op.Tags,
			}
			for k := range op.Req.Attrs.Content {
				item.Consumes = append(item.Consumes, k)
			}
			if _, ok := out[path]; !ok {
				out[path] = make(map[string]swaggerPathItem)
			}
			for _, v := range op.Tags {
				tags[v] = struct{}{}
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
