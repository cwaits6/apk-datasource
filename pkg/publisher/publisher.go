package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cwaits6/apk-datasource/pkg/generator"
)

// Publisher defines the interface for writing package data to a destination.
type Publisher interface {
	Publish(ctx context.Context, arch, packageName string, data *generator.RenovatePackage) error
}

// FilesystemPublisher writes JSON files to a local directory.
type FilesystemPublisher struct {
	BaseDir string
}

// NewFilesystemPublisher creates a new FilesystemPublisher.
func NewFilesystemPublisher(baseDir string) *FilesystemPublisher {
	return &FilesystemPublisher{BaseDir: baseDir}
}

// validatePackageName rejects names containing path traversal characters.
func validatePackageName(name string) error {
	if strings.Contains(name, "/") || strings.Contains(name, "..") || name == "" {
		return fmt.Errorf("invalid package name: %q", name)
	}
	return nil
}

// Publish writes a single package's Renovate JSON to disk atomically.
func (p *FilesystemPublisher) Publish(ctx context.Context, arch, packageName string, data *generator.RenovatePackage) error {
	if err := validatePackageName(packageName); err != nil {
		return err
	}
	if err := validatePackageName(arch); err != nil {
		return fmt.Errorf("invalid arch: %w", err)
	}

	dir := filepath.Join(p.BaseDir, arch)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON for %s/%s: %w", arch, packageName, err)
	}
	out = append(out, '\n')

	target := filepath.Join(dir, packageName+".json")
	tmp := target + ".tmp"

	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		return fmt.Errorf("writing temp file %s: %w", tmp, err)
	}

	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("renaming %s to %s: %w", tmp, target, err)
	}

	return nil
}
