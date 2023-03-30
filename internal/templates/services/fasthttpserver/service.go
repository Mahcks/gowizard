package fasthttpserver

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultPort            = "80"
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server *fasthttp.Server
	notify chan error
	opts   *Options
}

type Options struct {
	Port         string
	ReadTimeout  *time.Duration
	WriteTimeout *time.Duration
}

func New(handler fasthttp.RequestHandler, opts *Options) *Server {
	if opts.Port == "" {
		opts.Port = defaultPort
	}

	if opts.ReadTimeout == nil {
		opts.ReadTimeout = &defaultReadTimeout
	}

	if opts.WriteTimeout == nil {
		opts.WriteTimeout = &defaultWriteTimeout
	}

	httpServer := &fasthttp.Server{
		Handler:      handler,
		ReadTimeout:  *opts.ReadTimeout,
		WriteTimeout: *opts.WriteTimeout,
	}

	s := &Server{
		server: httpServer,
		notify: make(chan error, 1),
		opts:   opts,
	}

	s.start(fmt.Sprintf("0.0.0.0:%v", opts.Port))

	return s
}

func (s *Server) start(addr string) {
	go func() {
		s.notify <- s.server.ListenAndServe(addr)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown()
}
