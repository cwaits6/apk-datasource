package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cwaits6/apk-datasource/pkg/fetcher"
	"github.com/cwaits6/apk-datasource/pkg/generator"
	"github.com/rs/zerolog"
)

var knownArchitectures = map[string]bool{
	"x86_64":  true,
	"aarch64": true,
	"armhf":   true,
	"armv7":   true,
	"ppc64le": true,
	"s390x":   true,
	"riscv64": true,
}

// Server serves Renovate-compatible JSON over HTTP.
type Server struct {
	indexURLs       []string
	sourceURL       string
	homepage        string
	refreshInterval time.Duration
	port            int

	mu   sync.RWMutex
	data map[string]map[string]*generator.RenovatePackage // arch -> name -> package

	ready bool
}

// New creates a new Server instance.
func New(indexURLs []string, port int, refreshInterval time.Duration, sourceURL, homepage string) *Server {
	return &Server{
		indexURLs:       indexURLs,
		sourceURL:       sourceURL,
		homepage:        homepage,
		refreshInterval: refreshInterval,
		port:            port,
	}
}

// refresh fetches all indexes and updates the in-memory store.
func (s *Server) refresh(ctx context.Context) error {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("refreshing package data")

	sources, err := fetcher.FetchAll(ctx, s.indexURLs)
	if err != nil {
		return fmt.Errorf("fetching indexes: %w", err)
	}
	if len(sources) == 0 {
		return fmt.Errorf("no indexes fetched successfully")
	}

	data := generator.Generate(sources, s.sourceURL, s.homepage)

	s.mu.Lock()
	s.data = data
	s.ready = true
	s.mu.Unlock()

	total := 0
	for _, pkgs := range data {
		total += len(pkgs)
	}
	log.Info().Int("totalPackages", total).Int("architectures", len(data)).Msg("refresh complete")

	return nil
}

// Run starts the HTTP server with periodic refresh.
func (s *Server) Run(ctx context.Context) error {
	log := zerolog.Ctx(ctx)

	// Initial fetch (blocking).
	if err := s.refresh(ctx); err != nil {
		return fmt.Errorf("initial refresh: %w", err)
	}

	// Periodic refresh.
	go func() {
		ticker := time.NewTicker(s.refreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := s.refresh(ctx); err != nil {
					log.Error().Err(err).Msg("periodic refresh failed")
				}
			}
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.HandleFunc("GET /readyz", s.handleReadyz)
	mux.HandleFunc("GET /{arch}/{packageName}", s.handlePackage)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Graceful shutdown.
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Info().Msg("shutting down server")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("server shutdown error")
		}
	}()

	log.Info().Int("port", s.port).Msg("starting server")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

func (s *Server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handleReadyz(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	ready := s.ready
	s.mu.RUnlock()

	if !ready {
		http.Error(w, "not ready", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handlePackage(w http.ResponseWriter, r *http.Request) {
	arch := r.PathValue("arch")
	packageName := r.PathValue("packageName")

	if !knownArchitectures[arch] {
		http.Error(w, "unknown architecture", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	archData, ok := s.data[arch]
	if !ok {
		http.Error(w, "architecture not found", http.StatusNotFound)
		return
	}

	pkg, ok := archData[packageName]
	if !ok {
		http.Error(w, "package not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(pkg)
}
