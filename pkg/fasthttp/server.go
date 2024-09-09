package fasthttp

import (
	"fmt"
	"github.com/isaquesb/meli-url-shortener/internal/ports/input/http"
	"github.com/valyala/fasthttp"
)

type Server struct {
	Opts       http.Options
	Router     *Router
	FastServer *fasthttp.Server
}

func New(opts http.Options) http.Server {
	return &Server{
		Opts: opts,
	}
}

func (s *Server) Start(router http.Router) error {
	s.Router = router.(*Router)
	s.FastServer = &fasthttp.Server{
		Handler: s.Router.Handler,
	}
	return s.FastServer.ListenAndServe(fmt.Sprintf(":%d", s.Opts.Port))
}

func (s *Server) Shutdown() error {
	if s.FastServer == nil {
		return nil
	}
	return s.FastServer.Shutdown()
}

func (s *Server) Options() http.Options {
	return s.Opts
}
