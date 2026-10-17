package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Server struct {
	provider   externaldns.Provider
	httpServer *http.Server
}

func New(provider externaldns.Provider, listen string, readTimeout time.Duration, writeTimeout time.Duration) *Server {
	mux := http.NewServeMux()

	server := &Server{
		provider: provider,
		httpServer: &http.Server{
			Addr:              listen,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       readTimeout,
			WriteTimeout:      writeTimeout,
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

func (s *Server) Serve() error {
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go func(server *http.Server) {
		go func(server *http.Server) {
			defer waitGroup.Done()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			receivedSignal := <-sigChan

			shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), time.Minute)
			defer shutdownRelease()

			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Fatal().
					Err(err).
					Str("signal", receivedSignal.String()).
					Str("addr", server.Addr).
					Msg("HTTP shutdown error on")
			}

			log.Info().
				Str("signal", receivedSignal.String()).
				Str("addr", server.Addr).
				Msg("Graceful shutdown complete")
		}(server)

		log.Info().
			Str("addr", server.Addr).
			Msg("Starting server")

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().
				Err(err).
				Str("addr", server.Addr).
				Msg("Failed to run server")
		}
	}(s.httpServer)

	waitGroup.Wait()

	return nil
}
