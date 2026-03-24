package generator

import (
	"strings"

	"github.com/cwaits6/apk-datasource/pkg/fetcher"
	"github.com/cwaits6/apk-datasource/pkg/parser"
)

// Release represents a single version entry in the Renovate datasource output.
type Release struct {
	Version string `json:"version"`
}

// RenovatePackage is the Renovate custom datasource JSON schema for a single package.
type RenovatePackage struct {
	Releases  []Release `json:"releases"`
	SourceURL string    `json:"sourceUrl,omitempty"`
	Homepage  string    `json:"homepage,omitempty"`
}

// detectSourceURL infers a source repository URL from the APKINDEX URL.
func detectSourceURL(indexURL string) string {
	switch {
	case strings.Contains(indexURL, "packages.wolfi.dev"),
		strings.Contains(indexURL, "apk.cgr.dev/chainguard"):
		return "https://github.com/wolfi-dev/os"
	case strings.Contains(indexURL, "dl-cdn.alpinelinux.org"):
		return "https://gitlab.alpinelinux.org/alpine/aports"
	default:
		return ""
	}
}

// Generate transforms fetched index sources into Renovate-compatible JSON structures.
// Returns map[arch]map[packageName]*RenovatePackage.
func Generate(sources []*fetcher.IndexSource, sourceURLOverride, homepageOverride string) map[string]map[string]*RenovatePackage {
	result := make(map[string]map[string]*RenovatePackage)

	for _, src := range sources {
		parsed := parser.Parse(src.Packages)
		pkgMap := make(map[string]*RenovatePackage, len(parsed))

		sourceURL := sourceURLOverride
		if sourceURL == "" {
			sourceURL = detectSourceURL(src.URL)
		}

		for name, versions := range parsed {
			releases := make([]Release, 0, len(versions))
			for _, v := range versions {
				releases = append(releases, Release{Version: v.Version})
			}

			homepage := homepageOverride
			if homepage == "" && len(versions) > 0 {
				homepage = versions[0].URL
			}

			pkgMap[name] = &RenovatePackage{
				Releases:  releases,
				SourceURL: sourceURL,
				Homepage:  homepage,
			}
		}

		result[src.Arch] = pkgMap
	}

	return result
}
