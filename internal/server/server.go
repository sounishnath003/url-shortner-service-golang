package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sounishnath003/url-shortner-service-golang/internal/core"
	"github.com/sounishnath003/url-shortner-service-golang/internal/handlers"
	v1 "github.com/sounishnath003/url-shortner-service-golang/internal/handlers/v1"
	"golang.org/x/time/rate"
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

	// Auth endpoints.
	mux.HandleFunc("POST /login", handlers.LoginHandler)
	mux.HandleFunc("POST /signup", handlers.SignupHandler)

	// Groupping versioning.
	mux.HandleFunc("POST /api/v1/shorten", s.AuthGuardMiddleware(v1.GenerateUrlShortenerV1Handler))
	mux.HandleFunc("GET /api/v1/{shortenUrl}", v1.GetShortenUrlV1Handler)

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
	// ClientRateLimit helps to store the necessary information about the client.
	// Helps to idenitfy the client and last seen time. you can store the remote IP.
	type ClientRateLimit struct {
		limiter    *rate.Limiter
		remoteAddr string
		lastSeen   time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*ClientRateLimit)
	)

	// Run goroutine which will free and monitor the requests limit per client.
	go func() {
		// Run a infinite non blocking loop.
		for {
			// Add delay of running every minute.
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, client := range clients {
				// Allowing every client request 3 minute window
				// If the new request arrives more than 3 minute from the client.
				// We are not going to block / hence delete from the inmemory map
				// Note: can use redis here to make it distributed.
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO:
		// Extract remoteAddr.
		remoteAddr := r.RemoteAddr
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			// Throw 500 err.
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Lock the session for the client.
		mu.Lock()
		if _, found := clients[ip]; !found {
			// Create a new limiter for the client.
			clients[ip] = &ClientRateLimit{
				// rate.NewLimiter(2, 4), - means
				// 2 requests per second
				// 4 burst requests.
				limiter:    rate.NewLimiter(2, 4),
				remoteAddr: ip,
				lastSeen:   time.Now(),
			}
		}

		clients[ip].lastSeen = time.Now()
		// Check the rate limit for the client
		if !clients[ip].limiter.Allow() {
			// Unlock the session for the client.
			mu.Unlock()

			// Log the client information for audit trails
			s.co.Lo.Info("client", ip, "has been throttle due to too many requests", "lastSeen", clients[ip].lastSeen)

			w.Header().Add("X-Rate-Limiter", time.Now().Format(time.RFC3339Nano))
			w.Header().Add("Content-Type", "application/json, charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error":    "rate limit exceeded",
				"time":     time.Now().Format(time.RFC3339Nano),
				"ip":       ip,
				"msg":      "rate limit exceeded",
				"tryAfter": time.Now().Add(time.Second * 60).Format(time.RFC3339Nano),
			})
			return
		}

		mu.Unlock()

		w.Header().Add("X-Rate-Limiter", time.Now().Format(time.RFC3339Nano))
		next.ServeHTTP(w, r)
	}
}

// CustomReqContextMiddleware helps to feed the custom request context into the request context default channel.
// Which will inject the *core.Core dependencies execution in the handlers endpoints business logic.
func (s *Server) CustomReqContextMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Adding the *core.Core context into the default request context.
		ctx := context.WithValue(r.Context(), "co", s.co)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// HealthHandler works as a health check endpoint for the api.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	co := r.Context().Value("co").(*core.Core)

	hostname, err := os.Hostname()
	// Handle err.
	if err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	handlers.WriteJson(w, http.StatusOK, map[string]any{
		"status":    "OK",
		"version":   co.Version,
		"message":   "api services are normal",
		"hostname":  hostname,
		"timestamp": time.Now(),
	})
}

// AuthGuardMiddleware helps to authenticate the request.
// It checks the authorization header and verifies the JWT token.
// If the token is valid, the request is allowed to proceed.
// If the token is invalid, the request is rejected.
func (s *Server) AuthGuardMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.co.Lo.Info("inside auth middleware guard")

		authorization := r.Header.Get("Authorization")
		if len(authorization) < 5 {
			handlers.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}

		splits := strings.Split(authorization, " ")

		if splits[0] != "Bearer" || len(splits[1]) < 5 {
			handlers.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}

		token := splits[1]

		userID, foundEmail, err := s.ClaimAndVerifyJwtToken(token)
		if err != nil {
			s.co.Lo.Info("auth.middleware checks", "remoteIp", r.RemoteAddr, "isAuthorized", false)
			handlers.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		s.co.Lo.Info("auth.middleware checks", "remoteIp", r.RemoteAddr, "isAuthorized", true)
		// Set the request context with user
		ctx := context.WithValue(r.Context(), "userID", userID)
		ctx = context.WithValue(ctx, "userEmail", foundEmail)

		s.co.Lo.Info("checks completed auth middleware guard")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// ClaimAndVerifyJwtToken helps to verify the JWT token.
// It parses the token and checks if it is valid.
// If the token is valid, it returns the user's email address.
// If the token is invalid, it returns an error.
//
// Return userID, userEmail, error
func (s *Server) ClaimAndVerifyJwtToken(token string) (int, string, error) {
	// Parse the JWT token.
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.co.JwtSecret), nil
	})

	if err != nil {
		return 0, "", err
	}

	// Check if the token is valid.
	if !parsedToken.Valid {
		return 0, "", errors.New("Unauthorized")
	}

	// Check in database with the parsed token
	email, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return 0, "", errors.New("Unauthorized")
	}
	userID := 0
	foundEmail := ""
	foundPass := ""
	s.co.QueryStmts.GetUserByEmail.QueryRow(email).Scan(&userID, &foundEmail, &foundPass)
	if foundEmail == "" || foundPass == "" {
		return 0, "", errors.New("Unauthorized")
	}
	s.co.Lo.Info("user has been verified and authorized", "foundEmail", foundEmail)

	return userID, foundEmail, nil
}
