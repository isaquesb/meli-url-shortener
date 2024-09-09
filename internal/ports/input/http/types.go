package http

type Options struct {
	Port int
}

type RouteHandler func(Request) (Response, error)
