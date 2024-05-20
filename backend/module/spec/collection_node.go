package spec

import (
	"encoding/json"

	"github.com/apicat/apicat/v2/backend/module/spec/jsonschema"
)

type CollectionNode struct {
	Node
}

type Node interface {
	NodeType() string
}

type CollectionNodes []*CollectionNode

func NewCollectionNodesFromJson(c string) (CollectionNodes, error) {
	var collectionNodes CollectionNodes
	if err := json.Unmarshal([]byte(c), &collectionNodes); err != nil {
		return nil, err
	}
	return collectionNodes, nil
}

func (n *CollectionNode) ToHttpUrl() *CollectionHttpUrl {
	return n.Node.(*CollectionHttpUrl)
}

func (n *CollectionNode) ToHttpRequest() *CollectionHttpRequest {
	return n.Node.(*CollectionHttpRequest)
}

func (n *CollectionNode) ToHttpResponse() *CollectionHttpResponse {
	return n.Node.(*CollectionHttpResponse)
}

func (ns *CollectionNodes) DerefModel(ref *DefinitionModel) error {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().DerefModel(ref)
		case NODE_HTTP_RESPONSE:
			return node.ToHttpResponse().DerefModel(ref)
		}
	}
	return nil
}

func (ns *CollectionNodes) DeepDerefAll(params *GlobalParameters, definitions *Definitions) error {
	helper := jsonschema.NewDerefHelper(definitions.Schemas.ToJsonSchemaMap())

	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DerefGlobalParameters(params)
			if err := node.ToHttpRequest().DeepDerefModelByHelper(helper); err != nil {
				return err
			}
		case NODE_HTTP_RESPONSE:
			res := node.ToHttpResponse()
			if err := res.DerefAllResponses(definitions.Responses); err != nil {
				return err
			}
			if err := res.DeepDerefModelByHelper(helper); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ns *CollectionNodes) DerefGlobalParameter(in string, param *Parameter) {
	if param == nil {
		return
	}

	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().Attrs.Parameters.Add(in, param)
		}
	}
}

func (ns *CollectionNodes) DerefGlobalParameters(params *GlobalParameters) {
	if params == nil {
		return
	}

	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DerefGlobalParameters(params)
		}
	}
}

func (ns *CollectionNodes) DerefResponse(ref *DefinitionResponse) error {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			return node.ToHttpResponse().DerefResponse(ref)
		}
	}
	return nil
}

func (ns *CollectionNodes) DelRefModel(ref *DefinitionModel) {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DelRefModel(ref)
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DelRefModel(ref)
		}
	}
}

func (ns *CollectionNodes) DelRefResponse(ref *DefinitionResponse) {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().DelRefResponse(ref)
		}
	}
}

func (ns *CollectionNodes) DelGlobalExcept(in string, id int64) {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().DelGlobalExcept(in, id)
		}
	}
}

func (ns *CollectionNodes) GetGlobalExceptAll() map[string][]int64 {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			return node.ToHttpRequest().GetGlobalExceptAll()
		}
	}
	return nil
}

func (ns *CollectionNodes) GetRefModelIDs() []int64 {
	ids := make([]int64, 0)
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			ids = append(ids, node.ToHttpRequest().GetRefModelIDs()...)
		case NODE_HTTP_RESPONSE:
			ids = append(ids, node.ToHttpResponse().GetRefModelIDs()...)
		}
	}
	return ids
}

func (ns *CollectionNodes) GetRefResponseIDs() []int64 {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			return node.ToHttpResponse().GetRefResponseIDs()
		}
	}
	return nil
}

func (ns *CollectionNodes) AddReqParameter(in string, p *Parameter) {
	if p == nil {
		return
	}

	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_REQUEST:
			node.ToHttpRequest().Attrs.Parameters.Add(in, p)
		}
	}
}

func (ns *CollectionNodes) SortResponses() {
	for _, node := range *ns {
		switch node.NodeType() {
		case NODE_HTTP_RESPONSE:
			node.ToHttpResponse().Sort()
		}
	}
}

func (ns *CollectionNodes) GetUrlInfo() (method string, path string) {
	for _, node := range *ns {
		if node.NodeType() == NODE_HTTP_URL {
			url := node.ToHttpUrl()
			return url.Attrs.Method, url.Attrs.Path
		}
	}
	return method, path
}

func (ns *CollectionNodes) ToJson() (string, error) {
	res, err := json.Marshal(ns)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
