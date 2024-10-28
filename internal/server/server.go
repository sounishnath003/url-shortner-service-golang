package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

// Run intiaties or runs the backend http server.
// Initiaties all the handler endpoints on the specific port.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	// Adding the health endpoint.
	mux.HandleFunc("/healthy", HealthHandler)

	// Groupping versioning.
	mux.HandleFunc("POST /api/v1/shorten", GenerateUrlShortenerV1Handler)
	mux.HandleFunc("GET /api/v1/{shortenUrl}", GetShortenUrlV1Handler)

	hostAddr := fmt.Sprintf("http://0.0.0.0:%d", s.port)
	s.co.Lo.Info("server has been up and running", "on", hostAddr)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port),
		s.LoggerMiddleware(s.RateLimiterMiddleware(s.CustomReqContextMiddleware(mux))),
	)
}

// LoggerMiddleware helps to log every request received.
// Which helps for audit trails and service logs.
func (s *Server) LoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.co.Lo.Info("request received", "remoteAddr", r.RemoteAddr, "method", r.Method, "url", r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

// RateLimiterMiddleware rate limiter middleware to throttle ips which are brusting the traffic into the api server. It puts the detected Ips into a Minute pull back.
func (s *Server) RateLimiterMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		w.Header().Add("X-Rate-Limiter", time.Now().Format(time.RFC3339Nano))
		next.ServeHTTP(w, r)
	}
}

// CustomReqContextMiddleware helps to feed the custom request context into the request context default channel.
// Which will inject the *core.Core dependencies execution in the handlers endpoints business logic.
func (s *Server) CustomReqContextMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "co", s.co)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
