package spec

type HTTPParameters struct {
	Query  Schemas `json:"query"`
	Path   Schemas `json:"path"`
	Cookie Schemas `json:"cookie"`
	Header Schemas `json:"header"`
}

var HttpParameter = []string{"query", "path", "cookie", "header"}

func (h *HTTPParameters) Fill() {
	if h.Query == nil {
		h.Query = make(Schemas, 0)
	}
	if h.Path == nil {
		h.Path = make(Schemas, 0)
	}
	if h.Cookie == nil {
		h.Cookie = make(Schemas, 0)
	}
	if h.Header == nil {
		h.Header = make(Schemas, 0)
	}
}

func (h *HTTPParameters) Add(in string, v *Schema) {
	switch in {
	case "query":
		h.Query = append(h.Query, v)
	case "path":
		h.Path = append(h.Path, v)
	case "cookie":
		h.Cookie = append(h.Cookie, v)
	case "header":
		h.Header = append(h.Header, v)
	}
}

func (h *HTTPParameters) Map() map[string]Schemas {
	m := make(map[string]Schemas)
	if h.Query != nil {
		m["query"] = h.Query
	}
	if h.Path != nil {
		m["path"] = h.Path
	}
	if h.Header != nil {
		m["header"] = h.Header
	}
	if h.Cookie != nil {
		m["cookie"] = h.Cookie
	}
	return m
}

type HTTPNode[T HTTPNoder] struct {
	Type  string `json:"type"`
	Attrs T      `json:"attrs"`
}

type HTTPNoder interface {
	Name() string
}

func (n *HTTPNode[T]) NodeType() string {
	return n.Type
}

func WarpHTTPNode[T HTTPNoder](n T) Node {
	return &HTTPNode[T]{
		Type:  n.Name(),
		Attrs: n,
	}
}

type HTTPURLNode struct {
	Path   string  `json:"path"`
	Method string  `json:"method"`
	XDiff  *string `json:"x-apicat-diff,omitempty"`
}

func (HTTPURLNode) Name() string {
	return "apicat-http-url"
}

type HTTPBody map[string]*Schema

type HTTPRequestNode struct {
	GlobalExcepts map[string][]int64 `json:"globalExcepts,omitempty"`
	Parameters    HTTPParameters     `json:"parameters,omitempty"`
	Content       HTTPBody           `json:"content,omitempty"`
}

func (HTTPRequestNode) Name() string {
	return "apicat-http-request"
}

type HTTPResponsesNode struct {
	List HTTPResponses `json:"list,omitempty"`
}

func (HTTPResponsesNode) Name() string {
	return "apicat-http-response"
}

type HTTPResponse struct {
	Code  int     `json:"code"`
	XDiff *string `json:"x-apicat-diff,omitempty"`
	HTTPResponseDefine
}

type HTTPResponses []*HTTPResponse

func (h *HTTPResponses) Map() map[int]HTTPResponseDefine {
	m := make(map[int]HTTPResponseDefine)
	for _, v := range *h {
		m[v.Code] = v.HTTPResponseDefine
	}
	return m
}

func (h *HTTPResponses) Add(code int, hrd *HTTPResponseDefine) {
	for _, v := range *h {
		if v.Code == code {
			v.HTTPResponseDefine = *hrd
			v.XDiff = hrd.XDiff
			return
		}
	}
	*h = append(*h, &HTTPResponse{
		Code:               code,
		XDiff:              hrd.XDiff,
		HTTPResponseDefine: *hrd,
	})
}

func (h *HTTPResponse) SetXDiff(x *string) {
	h.Header.SetXDiff(x)
	h.Content.SetXDiff(x)
	h.HTTPResponseDefine.SetXDiff(x)
}

type HTTPResponseDefine struct {
	ID          int64    `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Content     HTTPBody `json:"content,omitempty"`
	Header      Schemas  `json:"header,omitempty"`
	Reference   *string  `json:"$ref,omitempty"`
	XDiff       *string  `json:"x-apicat-diff,omitempty"`
}

func (h *HTTPResponseDefine) Ref() bool { return h.Reference != nil }

func (h *HTTPResponseDefine) SetXDiff(x *string) {
	h.Header.SetXDiff(x)
	h.Content.SetXDiff(x)
	h.XDiff = x
}

type HTTPResponseDefines []HTTPResponseDefine

func (h HTTPResponseDefines) Lookup(name string) *HTTPResponseDefine {
	for _, v := range h {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func (h HTTPResponseDefines) LookupID(id int64) *HTTPResponseDefine {
	for _, v := range h {
		if v.ID == id {
			return &v
		}
	}
	return nil
}

func (h *HTTPResponseDefines) SetXDiff(x *string) {
	for _, v := range *h {
		v.SetXDiff(x)
	}
}

type HTTPPart struct {
	Title string
	ID    int64
	Dir   string
	HTTPRequestNode
	Responses HTTPResponses `json:"responses,omitempty"`
	XDiff     *string       `json:"x-apicat-diff,omitempty"`
}

func (h *HTTPPart) ToCollectItem(urlnode HTTPURLNode) *CollectItem {
	item := &CollectItem{
		Title: h.Title,
		Type:  ContentItemTypeHttp,
		XDiff: h.XDiff,
	}
	content := make([]*NodeProxy, 0)
	content = append(content, MuseCreateNodeProxy(WarpHTTPNode(urlnode)))
	content = append(content, MuseCreateNodeProxy(WarpHTTPNode(h.HTTPRequestNode)))
	content = append(content, MuseCreateNodeProxy(WarpHTTPNode(&HTTPResponsesNode{List: h.Responses})))
	item.Content = content
	return item
}

func (hb *HTTPBody) SetXDiff(x *string) {
	for _, v := range *hb {
		v.SetXDiff(x)
	}
}
