package fasthttp

import (
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/valyala/fasthttp"
)

func ParseResponse(res http.Response, ctx *fasthttp.RequestCtx) {
	ctx.SetBodyString(res.GetBody())
	for k, v := range res.GetHeaders() {
		ctx.Response.Header.Set(k, v)
	}
	ctx.SetStatusCode(res.GetStatusCode())
}
