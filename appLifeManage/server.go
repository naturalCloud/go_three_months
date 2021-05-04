package appLifeManage

import (
	"context"
	"net/http"
)

type Server interface {
	Start() error
	Stop() error
}

type HttpServer struct {
	server http.Server
	ctx    context.Context
}

func NewHttpServer(addr string, handler http.Handler) *HttpServer {
	return &HttpServer{server: http.Server{
		Addr:    addr,
		Handler: handler,
	}, ctx: context.Background()}
}

func (s *HttpServer) Start() error {
	err := s.server.ListenAndServe()

	return err
}

func (s *HttpServer) Stop() error {
	return s.server.Shutdown(s.ctx)
}

func NewHttpHandler() *http.ServeMux {
	return http.NewServeMux()
}
