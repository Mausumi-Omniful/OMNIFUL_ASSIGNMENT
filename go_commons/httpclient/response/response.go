package response

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Response interface {
	Body() []byte
	Status() string
	StatusCode() int
	Header() http.Header
	ContentType() string
	IsSuccess() bool
	IsError() bool
	IsCode2XX() bool
	UnmarshalBody(b any) error
	ForceFail()
}

func NewResponse(rc *resty.Client, rs *resty.Response) (Response, error) {
	return &response{
		Response:    rs,
		restyClient: rc,
	}, nil
}

func EmptyResponse() Response {
	return &response{
		Response: &resty.Response{},
	}
}

type response struct {
	*resty.Response

	restyClient *resty.Client
	forceFail   bool
}

func (r *response) ContentType() string {
	return r.Header().Get("Content-Type")
}

func (r *response) IsSuccess() bool {
	return !r.forceFail && r.Response.IsSuccess()
}

func (r *response) IsError() bool {
	return r.forceFail || r.Response.IsError()
}

func (r *response) IsCode2XX() bool {
	return r.StatusCode()/100 == 2
}

func (r *response) UnmarshalBody(b any) error {
	return resty.Unmarshalc(r.restyClient, r.ContentType(), r.Body(), &b)
}

func (r *response) ForceFail() {
	r.forceFail = true
}
