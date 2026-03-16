package parser

import (
	"sort"

	"gitlab.alpinelinux.org/alpine/go/repository"
)

// PackageInfo holds normalized metadata for a single package version.
type PackageInfo struct {
	Version string
	URL     string
	Origin  string
	License string
}

// Parse groups packages by name and returns a map of package name to sorted versions (descending).
func Parse(packages []*repository.Package) map[string][]PackageInfo {
	grouped := make(map[string][]PackageInfo)

	for _, pkg := range packages {
		if pkg.Name == "" || pkg.Version == "" {
			continue
		}

		info := PackageInfo{
			Version: pkg.Version,
			URL:     pkg.URL,
			Origin:  pkg.Origin,
			License: pkg.License,
		}

		grouped[pkg.Name] = append(grouped[pkg.Name], info)
	}

	// Deduplicate and sort versions descending by string comparison.
	for name, versions := range grouped {
		seen := make(map[string]bool)
		var deduped []PackageInfo
		for _, v := range versions {
			if !seen[v.Version] {
				seen[v.Version] = true
				deduped = append(deduped, v)
			}
		}
		sort.Slice(deduped, func(i, j int) bool {
			return deduped[i].Version > deduped[j].Version
		})
		grouped[name] = deduped
	}

	return grouped
}
