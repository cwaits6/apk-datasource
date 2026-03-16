//go:build ignore

// This program generates the testdata/APKINDEX.tar.gz fixture using the Alpine library.
// Run with: go run internal/testdata/generate_fixture.go
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"gitlab.alpinelinux.org/alpine/go/repository"
)

func main() {
	packages := []*repository.Package{
		{
			Name:        "curl",
			Version:     "8.11.1-r0",
			Arch:        "x86_64",
			Description: "URL retrieval utility and library",
			URL:         "https://curl.se",
			License:     "MIT",
			Origin:      "curl",
			Maintainer:  "Wolfi <wolfi@chainguard.dev>",
			Checksum:    []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14},
			Size:        300000,
			BuildTime:   time.Unix(1700000000, 0).UTC(),
		},
		{
			Name:        "git",
			Version:     "2.43.0-r0",
			Arch:        "x86_64",
			Description: "Distributed version control system",
			URL:         "https://git-scm.com",
			License:     "GPL-2.0-only",
			Origin:      "git",
			Maintainer:  "Wolfi <wolfi@chainguard.dev>",
			Checksum:    []byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23, 0x24},
			Size:        500000,
			BuildTime:   time.Unix(1700000000, 0).UTC(),
		},
		{
			Name:        "openssl",
			Version:     "3.2.0-r0",
			Arch:        "x86_64",
			Description: "Toolkit for Transport Layer Security (TLS)",
			URL:         "https://www.openssl.org",
			License:     "Apache-2.0",
			Origin:      "openssl",
			Maintainer:  "Wolfi <wolfi@chainguard.dev>",
			Checksum:    []byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34},
			Size:        400000,
			BuildTime:   time.Unix(1700000000, 0).UTC(),
		},
		{
			Name:        "c++",
			Version:     "13.2.0-r0",
			Arch:        "x86_64",
			Description: "C++ compiler",
			URL:         "https://gcc.gnu.org",
			License:     "GPL-3.0-or-later",
			Origin:      "gcc",
			Maintainer:  "Wolfi <wolfi@chainguard.dev>",
			Checksum:    []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40, 0x41, 0x42, 0x43, 0x44},
			Size:        200000,
			BuildTime:   time.Unix(1700000000, 0).UTC(),
		},
	}

	idx := &repository.ApkIndex{
		Packages: packages,
	}

	archive, err := repository.ArchiveFromIndex(idx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ArchiveFromIndex: %v\n", err)
		os.Exit(1)
	}

	outDir := filepath.Join("internal", "testdata")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}

	outPath := filepath.Join(outDir, "APKINDEX.tar.gz")
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	n, err := io.Copy(f, archive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "copy: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s (%d bytes) with %d packages\n", outPath, n, len(packages))
}
