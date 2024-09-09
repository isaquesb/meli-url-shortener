package http

import "context"

type Request interface {
	Ctx() context.Context
	IsGet() bool
	IsPost() bool
	IsDelete() bool
	FormValue(key string) string
	PathValue(key string) any
	Header(key string) string
}

type Router interface {
	GET(path string, handler RouteHandler)
	POST(path string, handler RouteHandler)
	DELETE(path string, handler RouteHandler)
}

type Response interface {
	Header(key string) string
	GetBody() string
	GetStatusCode() int
	GetHeaders() map[string]string
}

type Server interface {
	Start(Router) error
	Shutdown() error
	Options() Options
}
