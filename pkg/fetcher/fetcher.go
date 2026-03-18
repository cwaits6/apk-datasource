package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/rs/zerolog"
	"gitlab.alpinelinux.org/alpine/go/repository"
	"golang.org/x/sync/errgroup"
)

var archPattern = regexp.MustCompile(`/(x86_64|aarch64|armhf|armv7|ppc64le|s390x|riscv64)/`)

// IndexSource holds the parsed packages from a single APKINDEX along with metadata.
type IndexSource struct {
	Arch     string
	Packages []*repository.Package
	URL      string
}

// detectArch extracts the architecture from an APKINDEX URL path.
func detectArch(url string) string {
	matches := archPattern.FindStringSubmatch(url)
	if len(matches) >= 2 {
		return matches[1]
	}
	return "unknown"
}

// Fetch downloads and parses a single APKINDEX.tar.gz URL.
func Fetch(ctx context.Context, url string) (*IndexSource, error) {
	log := zerolog.Ctx(ctx)

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request for %s: %w", url, err)
	}

	log.Info().Str("url", url).Msg("fetching APKINDEX")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil, fmt.Errorf("fetching %s: HTTP %d", url, resp.StatusCode)
	}

	idx, err := repository.IndexFromArchive(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing index from %s: %w", url, err)
	}

	arch := detectArch(url)
	log.Info().Str("arch", arch).Int("packages", len(idx.Packages)).Msg("parsed APKINDEX")

	return &IndexSource{
		Arch:     arch,
		Packages: idx.Packages,
		URL:      url,
	}, nil
}

// FetchAll downloads and parses multiple APKINDEX URLs concurrently.
func FetchAll(ctx context.Context, urls []string) ([]*IndexSource, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(4)

	results := make([]*IndexSource, len(urls))

	for i, u := range urls {
		g.Go(func() error {
			src, err := Fetch(ctx, u)
			if err != nil {
				zerolog.Ctx(ctx).Warn().Err(err).Str("url", u).Msg("failed to fetch index, skipping")
				return nil // partial failure: continue
			}
			results[i] = src
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Filter nil results from failed fetches.
	var filtered []*IndexSource
	for _, r := range results {
		if r != nil {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}
