package parser

import (
	"testing"

	"gitlab.alpinelinux.org/alpine/go/repository"
)

func TestParse_Grouping(t *testing.T) {
	packages := []*repository.Package{
		{Name: "curl", Version: "8.11.1-r0", URL: "https://curl.se", License: "MIT", Origin: "curl"},
		{Name: "git", Version: "2.43.0-r0", URL: "https://git-scm.com", License: "GPL-2.0-only", Origin: "git"},
		{Name: "curl", Version: "8.10.0-r0", URL: "https://curl.se", License: "MIT", Origin: "curl"},
	}

	result := Parse(packages)

	if len(result) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(result))
	}

	curlVersions := result["curl"]
	if len(curlVersions) != 2 {
		t.Fatalf("expected 2 curl versions, got %d", len(curlVersions))
	}

	// Should be sorted descending.
	if curlVersions[0].Version != "8.11.1-r0" {
		t.Errorf("expected first version 8.11.1-r0, got %s", curlVersions[0].Version)
	}
	if curlVersions[1].Version != "8.10.0-r0" {
		t.Errorf("expected second version 8.10.0-r0, got %s", curlVersions[1].Version)
	}
}

func TestParse_Deduplication(t *testing.T) {
	packages := []*repository.Package{
		{Name: "curl", Version: "8.11.1-r0", URL: "https://curl.se"},
		{Name: "curl", Version: "8.11.1-r0", URL: "https://curl.se"},
	}

	result := Parse(packages)
	if len(result["curl"]) != 1 {
		t.Errorf("expected 1 version after dedup, got %d", len(result["curl"]))
	}
}

func TestParse_EmptyInput(t *testing.T) {
	result := Parse(nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil input, got %d", len(result))
	}
}

func TestParse_SkipsMissingFields(t *testing.T) {
	packages := []*repository.Package{
		{Name: "", Version: "1.0.0"},
		{Name: "valid", Version: ""},
		{Name: "good", Version: "1.0.0"},
	}

	result := Parse(packages)
	if len(result) != 1 {
		t.Errorf("expected 1 package, got %d", len(result))
	}
	if _, ok := result["good"]; !ok {
		t.Error("expected 'good' package")
	}
}

func TestParse_PreservesMetadata(t *testing.T) {
	packages := []*repository.Package{
		{Name: "curl", Version: "8.11.1-r0", URL: "https://curl.se", License: "MIT", Origin: "curl"},
	}

	result := Parse(packages)
	info := result["curl"][0]

	if info.URL != "https://curl.se" {
		t.Errorf("expected URL https://curl.se, got %s", info.URL)
	}
	if info.License != "MIT" {
		t.Errorf("expected license MIT, got %s", info.License)
	}
	if info.Origin != "curl" {
		t.Errorf("expected origin curl, got %s", info.Origin)
	}
}
