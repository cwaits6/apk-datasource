package publisher

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/cwaits6/apk-datasource/pkg/generator"
)

func TestFilesystemPublisher_Write(t *testing.T) {
	dir := t.TempDir()
	pub := NewFilesystemPublisher(dir)

	data := &generator.RenovatePackage{
		Releases: []generator.Release{
			{Version: "8.11.1-r0"},
		},
		SourceURL: "https://github.com/wolfi-dev/os",
		Homepage:  "https://curl.se",
	}

	if err := pub.Publish(context.Background(), "x86_64", "curl", data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outPath := filepath.Join(dir, "x86_64", "curl.json")
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	var result generator.RenovatePackage
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(result.Releases) != 1 || result.Releases[0].Version != "8.11.1-r0" {
		t.Errorf("unexpected content: %+v", result)
	}
}

func TestFilesystemPublisher_PathTraversal(t *testing.T) {
	dir := t.TempDir()
	pub := NewFilesystemPublisher(dir)

	tests := []struct {
		name string
		pkg  string
	}{
		{"dot-dot", "../evil"},
		{"slash", "foo/bar"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pub.Publish(context.Background(), "x86_64", tt.pkg, &generator.RenovatePackage{})
			if err == nil {
				t.Error("expected error for invalid package name")
			}
		})
	}
}

func TestFilesystemPublisher_InvalidArch(t *testing.T) {
	dir := t.TempDir()
	pub := NewFilesystemPublisher(dir)

	err := pub.Publish(context.Background(), "../evil", "curl", &generator.RenovatePackage{})
	if err == nil {
		t.Error("expected error for invalid arch")
	}
}

func TestFilesystemPublisher_SpecialCharInName(t *testing.T) {
	dir := t.TempDir()
	pub := NewFilesystemPublisher(dir)

	data := &generator.RenovatePackage{
		Releases: []generator.Release{{Version: "13.2.0-r0"}},
	}

	// Package name with + (like c++)
	if err := pub.Publish(context.Background(), "x86_64", "c++", data); err != nil {
		t.Fatalf("unexpected error for c++: %v", err)
	}

	outPath := filepath.Join(dir, "x86_64", "c++.json")
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("expected c++.json to exist")
	}
}
