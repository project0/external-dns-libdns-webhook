package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	provider     externaldns.Provider
	httpServer   *http.Server
	healthServer *http.Server
}

// Config of the webhook servers.
type Config struct {
	// Listen address of the webhook API, external-dns expects it on localhost:8888 by default.
	Listen string
	// HealthListen address serving /healthz and /readyz for probes, disabled if empty.
	HealthListen string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

const readHeaderTimeout = 5 * time.Second

func New(provider externaldns.Provider, config Config) *Server {
	mux := http.NewServeMux()

	server := &Server{
		provider: provider,
		httpServer: &http.Server{
			Addr:              config.Listen,
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       config.ReadTimeout,
			WriteTimeout:      config.WriteTimeout,
			Handler:           loggerMiddleware(mux),
		},
	}

	// method patterns respond with 405 on other methods, "/{$}" only matches the root path
	mux.HandleFunc("GET /{$}", server.NegotiateHandler)
	mux.HandleFunc("GET /records", server.RecordsHandler)
	mux.HandleFunc("POST /records", server.ApplyChangesHandler)
	mux.HandleFunc("POST /adjustendpoints", server.AdjustEndpointsHandler)

	mux.HandleFunc("GET /healthz", okHandler)
	mux.HandleFunc("GET /readyz", okHandler)

	if config.HealthListen != "" {
		healthMux := http.NewServeMux()
		healthMux.HandleFunc("GET /healthz", okHandler)
		healthMux.HandleFunc("GET /readyz", okHandler)

		server.healthServer = &http.Server{
			Addr:              config.HealthListen,
			ReadHeaderTimeout: readHeaderTimeout,
			ReadTimeout:       config.ReadTimeout,
			WriteTimeout:      config.WriteTimeout,
			Handler:           healthMux,
		}
	}

	return server
}

func okHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "ok", http.StatusOK)
}

// Handler returns the HTTP handler of the webhook API.
func (s *Server) Handler() http.Handler {
	return s.httpServer.Handler
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.With().
			Str("request_uri", r.RequestURI).
			Str("method", r.Method).
			Logger()
		r = r.WithContext(logger.WithContext(r.Context()))

		// only intercept bodies if trace logging is enabled
		if zerolog.GlobalLevel() > zerolog.TraceLevel {
			next.ServeHTTP(w, r)

			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Err(err).Msg("Failed to read request body")
		}

		logger.Trace().RawJSON("body", jsonOrNull(body)).Msg("Request body")
		r.Body = io.NopCloser(bytes.NewReader(body))

		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		logger.Trace().Int("status", rec.status).RawJSON("body", jsonOrNull(rec.body.Bytes())).Msg("Response body")
	})
}

func jsonOrNull(body []byte) []byte {
	if !json.Valid(body) {
		quoted, _ := json.Marshal(string(body)) //nolint:errchkjson // string is always encodable

		return quoted
	}

	return body
}

type responseRecorder struct {
	http.ResponseWriter

	status int
	body   bytes.Buffer
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)

	n, err := r.ResponseWriter.Write(b)
	if err != nil {
		return n, fmt.Errorf("failed to write response: %w", err)
	}

	return n, nil
}

// Serve runs the servers until SIGINT or SIGTERM is received.
func (s *Server) Serve() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	servers := []*http.Server{s.httpServer}
	if s.healthServer != nil {
		servers = append(servers, s.healthServer)
	}

	errs := make(chan error, len(servers))

	for _, server := range servers {
		go func() {
			log.Info().Str("addr", server.Addr).Msg("Starting server")

			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errs <- fmt.Errorf("server %s: %w", server.Addr, err)
			}
		}()
	}

	var serveErr error

	select {
	case <-ctx.Done():
		log.Info().Msg("Received shutdown signal")
	case serveErr = <-errs:
	}

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), time.Minute)
	defer shutdownRelease()

	for _, server := range servers {
		if err := server.Shutdown(shutdownCtx); err != nil {
			serveErr = errors.Join(serveErr, fmt.Errorf("shutdown server %s: %w", server.Addr, err))
		}
	}

	if serveErr == nil {
		log.Info().Msg("Graceful shutdown complete")
	}

	return serveErr
}
