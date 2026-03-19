package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cwaits6/apk-datasource/pkg/generator"
	"github.com/cwaits6/apk-datasource/pkg/metrics"
)

func newTestServer(t *testing.T) (*Server, *httptest.Server) {
	t.Helper()

	data, err := os.ReadFile("../../internal/testdata/APKINDEX.tar.gz")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write(data); err != nil {
			t.Errorf("failed to write fixture data: %v", err)
		}
	}))

	srv := New(
		[]string{upstream.URL + "/x86_64/APKINDEX.tar.gz"},
		0, // port unused for handler tests
		1*time.Hour,
		"", "",
		metrics.Noop(), nil,
	)

	ctx := context.Background()
	if err := srv.refresh(ctx); err != nil {
		t.Fatalf("refresh failed: %v", err)
	}

	return srv, upstream
}

func TestServer_HandlePackage(t *testing.T) {
	srv, upstream := newTestServer(t)
	defer upstream.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{arch}/{packageName}", srv.handlePackage)

	req := httptest.NewRequest("GET", "/x86_64/curl", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var pkg generator.RenovatePackage
	if err := json.NewDecoder(w.Body).Decode(&pkg); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(pkg.Releases) == 0 {
		t.Error("expected releases")
	}
}

func TestServer_HandlePackage_NotFound(t *testing.T) {
	srv, upstream := newTestServer(t)
	defer upstream.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{arch}/{packageName}", srv.handlePackage)

	req := httptest.NewRequest("GET", "/x86_64/nonexistent-package-xyz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestServer_HandlePackage_UnknownArch(t *testing.T) {
	srv, upstream := newTestServer(t)
	defer upstream.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{arch}/{packageName}", srv.handlePackage)

	req := httptest.NewRequest("GET", "/mips64/curl", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestServer_Healthz(t *testing.T) {
	srv, upstream := newTestServer(t)
	defer upstream.Close()

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	srv.handleHealthz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestServer_Readyz_BeforeLoad(t *testing.T) {
	srv := &Server{}

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()
	srv.handleReadyz(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503 before load, got %d", w.Code)
	}
}

func TestServer_Metrics(t *testing.T) {
	m, handler, err := metrics.Setup()
	if err != nil {
		t.Fatalf("metrics setup: %v", err)
	}

	srv := &Server{
		metrics:        m,
		metricsHandler: handler,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", srv.handleHealthz)
	mux.Handle("GET /metrics", srv.metricsHandler)

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		t.Errorf("expected text/plain content type, got %s", contentType)
	}
}

func TestServer_Readyz_AfterLoad(t *testing.T) {
	srv, upstream := newTestServer(t)
	defer upstream.Close()

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()
	srv.handleReadyz(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 after load, got %d", w.Code)
	}
}
