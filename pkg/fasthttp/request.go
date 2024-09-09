package fasthttp

import (
	"context"
	"github.com/valyala/fasthttp"
)

type Request struct {
	Context *fasthttp.RequestCtx
}

func NewRequest(ctx *fasthttp.RequestCtx) *Request {
	return &Request{
		Context: ctx,
	}
}

func (r *Request) Ctx() context.Context {
	return r.Context
}

func (r *Request) IsGet() bool {
	return r.Context.IsGet()
}

func (r *Request) IsPost() bool {
	return r.Context.IsPost()
}

func (r *Request) IsDelete() bool {
	return r.Context.IsDelete()
}

func (r *Request) FormValue(key string) string {
	return string(r.Context.FormValue(key))
}

func (r *Request) PathValue(key string) any {
	return r.Context.UserValue(key)
}

func (r *Request) Header(key string) string {
	return string(r.Context.Request.Header.Peek(key))
}
