package request

import (
	"github.com/omniful/go_commons/util"
	"net/http"
	"net/url"
)

var defaultMethod = http.MethodGet

type Builder interface {
	SetMethod(string) Builder
	SetUri(string) Builder
	SetHeaders(url.Values) Builder
	SetPathParams(PathParams) Builder
	SetQueryParams(url.Values) Builder
	SetBody(any) Builder
	SetFormData(values url.Values) Builder
	Build() (Request, error)
}

func NewBuilder() Builder {
	return &builder{
		h:  make(url.Values),
		pp: make(PathParams),
		qp: make(url.Values),
		fd: make(url.Values),
	}
}

type builder struct {
	method string
	uri    string
	h      url.Values
	pp     PathParams
	qp     url.Values
	body   any
	fd     url.Values
}

func (b *builder) SetMethod(m string) Builder {
	b.method = m
	return b
}

func (b *builder) SetUri(u string) Builder {
	b.uri = u
	return b
}

func (b *builder) SetHeaders(h url.Values) Builder {
	b.h = h
	return b
}

func (b *builder) SetPathParams(pp PathParams) Builder {
	b.pp = pp
	return b
}

func (b *builder) SetQueryParams(qp url.Values) Builder {
	b.qp = qp
	return b
}

func (b *builder) SetBody(bd any) Builder {
	b.body = bd
	return b
}

func (b *builder) SetFormData(fd url.Values) Builder {
	b.fd = fd
	return b
}

func (b *builder) Build() (Request, error) {
	m := util.FirstNonEmptyString(b.method, defaultMethod)
	return &request{
		Method:      m,
		Uri:         b.uri,
		Headers:     b.h,
		PathParams:  b.pp,
		QueryParams: b.qp,
		Body:        b.body,
		FormData:    b.fd,
	}, nil
}
