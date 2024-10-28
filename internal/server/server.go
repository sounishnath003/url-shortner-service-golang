package server

import (
	"fmt"
	"net/http"

	"github.com/sounishnath003/url-shortner-service-golang/internal/core"
)

type Server struct {
	port int
	co   *core.Core
}

func NewServer(co *core.Core) *Server {
	return &Server{
		port: co.Port,
		co:   co,
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()

	// Adding the health endpoint.
	mux.HandleFunc("/healthy", HealthHandler)

	// Groupping versioning.
	mux.HandleFunc("POST /api/v1/shorten", GenerateUrlShortenerV1Handler)
	mux.HandleFunc("GET /api/v1/{shortenUrl}", GetShortenUrlV1Handler)

	hostAddr := fmt.Sprintf("http://0.0.0.0:%d", s.port)
	s.co.Lo.Info("server has been up and running", "on", hostAddr)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.LoggerMiddleware(mux))
}

// LoggerMiddleware helps to log every request received, which helps for audit trails and service logs.
func (s *Server) LoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.co.Lo.Info("request received", "remoteAddr", r.RemoteAddr, "method", r.Method, "url", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}
