package request

import "net/url"

type PathParams map[string]string

type Request interface {
	GetMethod() string
	GetUri() string
	GetHeaders() url.Values
	GetPathParams() PathParams
	GetQueryParams() url.Values
	GetBody() any
	GetFormData() url.Values
	ToBuilder() Builder
}

type request struct {
	Method      string
	Uri         string
	Headers     url.Values
	PathParams  PathParams
	QueryParams url.Values
	Body        any
	FormData    url.Values
}

func (r *request) GetMethod() string          { return r.Method }
func (r *request) GetUri() string             { return r.Uri }
func (r *request) GetHeaders() url.Values     { return r.Headers }
func (r *request) GetPathParams() PathParams  { return r.PathParams }
func (r *request) GetQueryParams() url.Values { return r.QueryParams }
func (r *request) GetBody() any               { return r.Body }
func (r *request) GetFormData() url.Values    { return r.FormData }

func (r *request) ToBuilder() Builder {
	return NewBuilder().
		SetMethod(r.Method).
		SetUri(r.Uri).
		SetHeaders(r.Headers).
		SetPathParams(r.PathParams).
		SetQueryParams(r.QueryParams).
		SetBody(r.Body).
		SetFormData(r.FormData)
}
