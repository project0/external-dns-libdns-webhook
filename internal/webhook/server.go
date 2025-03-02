package webhook

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
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

	mux.HandleFunc("/", server.NegotiateHandler)
	mux.HandleFunc("/records", server.RecordsHandler)
	mux.HandleFunc("/adjustendpoints", server.AdjustendpointsHandler)

	mux.HandleFunc("/healthz", okHandler)
	mux.HandleFunc("/readyz", okHandler)

	return server
}

func okHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "ok", http.StatusOK)
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := log.With().
			Str("request_uri", r.RequestURI).
			Str("method", r.Method).
			Logger().
			WithContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
