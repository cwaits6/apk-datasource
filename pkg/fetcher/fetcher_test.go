package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func fixtureServer(t *testing.T) *httptest.Server {
	t.Helper()
	data, err := os.ReadFile("../../internal/testdata/APKINDEX.tar.gz")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(data)
	}))
}

func TestFetch_Success(t *testing.T) {
	srv := fixtureServer(t)
	defer srv.Close()

	url := srv.URL + "/x86_64/APKINDEX.tar.gz"
	src, err := Fetch(context.Background(), url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if src.Arch != "x86_64" {
		t.Errorf("expected arch x86_64, got %s", src.Arch)
	}
	if len(src.Packages) == 0 {
		t.Error("expected packages, got none")
	}
	if src.URL != url {
		t.Errorf("expected URL %s, got %s", url, src.URL)
	}
}

func TestFetch_404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	_, err := Fetch(context.Background(), srv.URL+"/x86_64/APKINDEX.tar.gz")
	if err == nil {
		t.Fatal("expected error for 404")
	}
}

func TestFetch_500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := Fetch(context.Background(), srv.URL+"/x86_64/APKINDEX.tar.gz")
	if err == nil {
		t.Fatal("expected error for 500")
	}
}

func TestFetch_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := Fetch(ctx, srv.URL+"/x86_64/APKINDEX.tar.gz")
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestDetectArch(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz", "x86_64"},
		{"https://packages.wolfi.dev/os/aarch64/APKINDEX.tar.gz", "aarch64"},
		{"https://dl-cdn.alpinelinux.org/alpine/v3.19/main/armhf/APKINDEX.tar.gz", "armhf"},
		{"https://example.com/APKINDEX.tar.gz", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := detectArch(tt.url)
			if got != tt.expected {
				t.Errorf("detectArch(%q) = %q, want %q", tt.url, got, tt.expected)
			}
		})
	}
}

func TestFetchAll_PartialFailure(t *testing.T) {
	srv := fixtureServer(t)
	defer srv.Close()

	badSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "fail", http.StatusInternalServerError)
	}))
	defer badSrv.Close()

	urls := []string{
		srv.URL + "/x86_64/APKINDEX.tar.gz",
		badSrv.URL + "/aarch64/APKINDEX.tar.gz",
	}

	results, err := FetchAll(context.Background(), urls)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should get 1 success (bad URL is skipped).
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}
