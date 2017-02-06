// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package graceful

import (
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	middleware "github.com/justinas/alice"
	"gopkg.in/tylerb/graceful.v1"
)

const timeout = 5 * time.Second

type HTTPServer struct {
	sync.Mutex

	network string
	router  *http.ServeMux
	server  *graceful.Server
	err     error
}

func recovery(handler http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				http.Error(w, "internal error, check logs for details", http.StatusInternalServerError)
			}
		}()
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func NewHTTPServer(net, addr string, mw ...middleware.Constructor) *HTTPServer {
	r := http.NewServeMux()

	return &HTTPServer{
		network: net,
		router:  r,
		server: &graceful.Server{
			Timeout: timeout,
			Server: &http.Server{
				Addr:         addr,
				Handler:      middleware.New(recovery).Append(mw...).Then(r),
				ReadTimeout:  timeout,
				WriteTimeout: timeout,
			},
		},
	}
}

func (s *HTTPServer) Handle(method, route string, handler http.HandlerFunc) {
	f := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.NotFound(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	}
	s.router.HandleFunc(route, f)
}

func (s *HTTPServer) Serve() <-chan struct{} {
	if s.network == "unix" {
		os.Remove(s.server.Addr)
	}
	l, err := net.Listen(s.network, s.server.Addr)
	if err != nil {
		s.Lock()
		s.err = err
		s.Unlock()
		c := make(chan struct{})
		close(c)
		return c
	}

	c := s.server.StopChan()
	go func() {
		s.Lock()
		defer s.Unlock()

		err = s.server.Serve(l)
		if e, ok := err.(*net.OpError); !ok || (ok && e.Op != "accept") {
			s.err = err
		}
	}()
	return c
}

func (s *HTTPServer) Stop() {
	s.server.Stop(timeout)
}

func (s *HTTPServer) Error() error {
	s.Lock()
	defer s.Unlock()

	return s.err
}
