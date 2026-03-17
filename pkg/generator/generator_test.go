package generator

import (
	"encoding/json"
	"testing"

	"github.com/cwaits6/apk-datasource/pkg/fetcher"
	"gitlab.alpinelinux.org/alpine/go/repository"
)

func TestGenerate_BasicOutput(t *testing.T) {
	sources := []*fetcher.IndexSource{
		{
			Arch: "x86_64",
			URL:  "https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz",
			Packages: []*repository.Package{
				{Name: "curl", Version: "8.11.1-r0", URL: "https://curl.se"},
				{Name: "git", Version: "2.43.0-r0", URL: "https://git-scm.com"},
			},
		},
	}

	result := Generate(sources, "", "")

	if len(result) != 1 {
		t.Fatalf("expected 1 arch, got %d", len(result))
	}

	x86 := result["x86_64"]
	if len(x86) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(x86))
	}

	curl := x86["curl"]
	if curl == nil {
		t.Fatal("expected curl package")
	}
	if len(curl.Releases) != 1 {
		t.Fatalf("expected 1 release, got %d", len(curl.Releases))
	}
	if curl.Releases[0].Version != "8.11.1-r0" {
		t.Errorf("expected version 8.11.1-r0, got %s", curl.Releases[0].Version)
	}
}

func TestGenerate_SourceURLDetection(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{"wolfi", "https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz", "https://github.com/wolfi-dev/os"},
		{"alpine", "https://dl-cdn.alpinelinux.org/alpine/v3.19/main/x86_64/APKINDEX.tar.gz", "https://gitlab.alpinelinux.org/alpine/aports"},
		{"unknown", "https://example.com/APKINDEX.tar.gz", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectSourceURL(tt.url)
			if got != tt.expected {
				t.Errorf("detectSourceURL(%q) = %q, want %q", tt.url, got, tt.expected)
			}
		})
	}
}

func TestGenerate_OverrideSourceURL(t *testing.T) {
	sources := []*fetcher.IndexSource{
		{
			Arch: "x86_64",
			URL:  "https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz",
			Packages: []*repository.Package{
				{Name: "curl", Version: "8.11.1-r0"},
			},
		},
	}

	result := Generate(sources, "https://custom.example.com", "https://homepage.example.com")
	curl := result["x86_64"]["curl"]

	if curl.SourceURL != "https://custom.example.com" {
		t.Errorf("expected custom source URL, got %s", curl.SourceURL)
	}
	if curl.Homepage != "https://homepage.example.com" {
		t.Errorf("expected custom homepage, got %s", curl.Homepage)
	}
}

func TestGenerate_RenovateJSONContract(t *testing.T) {
	sources := []*fetcher.IndexSource{
		{
			Arch: "x86_64",
			URL:  "https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz",
			Packages: []*repository.Package{
				{Name: "curl", Version: "8.11.1-r0", URL: "https://curl.se"},
			},
		},
	}

	result := Generate(sources, "", "")
	curl := result["x86_64"]["curl"]

	// Marshal and unmarshal to verify JSON structure.
	data, err := json.Marshal(curl)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var roundTrip RenovatePackage
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if len(roundTrip.Releases) != 1 {
		t.Fatalf("expected 1 release after round-trip, got %d", len(roundTrip.Releases))
	}
	if roundTrip.Releases[0].Version != "8.11.1-r0" {
		t.Errorf("expected version 8.11.1-r0 after round-trip, got %s", roundTrip.Releases[0].Version)
	}

	// Verify JSON has correct keys.
	var raw map[string]interface{}
	_ = json.Unmarshal(data, &raw)
	if _, ok := raw["releases"]; !ok {
		t.Error("JSON missing 'releases' key")
	}
}
