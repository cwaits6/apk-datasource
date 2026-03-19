# CLAUDE.md

## Project Overview

apk-datasource generates Renovate-compatible custom datasource JSON from Wolfi and Alpine APK package indexes. It solves the problem of auto-updating pinned `apk add package=version` lines in Dockerfiles (see renovatebot/renovate#5422).

Two modes: `generate` (static files to disk) and `serve` (HTTP server with periodic refresh).

## Architecture

```
cmd/apk-datasource/          CLI entry point (cobra)
  cmd/root.go                 Root command, logging setup (zerolog)
  cmd/serve.go                HTTP server command with --metrics flag
  cmd/generate.go             Static file generation command
  cmd/version.go              Version info (set via ldflags)
pkg/fetcher/                  Downloads and parses APKINDEX.tar.gz (concurrent, max 4)
pkg/parser/                   Groups packages by name, deduplicates versions
pkg/generator/                Converts to Renovate JSON schema (RenovatePackage)
pkg/publisher/                Writes JSON files atomically (filesystem publisher)
pkg/server/                   HTTP server with in-memory store, RWMutex, periodic refresh
pkg/metrics/                  OTel-Prometheus bridge (Setup/Noop pattern)
```

## Key Design Decisions

- **Single port**: `/metrics`, `/healthz`, `/readyz`, and `/{arch}/{packageName}` all on the same mux
- **Nil-safe Metrics**: `Noop()` returns `&Metrics{}` with nil instruments; all methods guard on `m == nil`
- **OTel-Prometheus bridge**: Single setup, OTel-native code, Prometheus-compatible `/metrics` output
- **Path normalization**: Metrics middleware maps paths to route templates (`/{arch}/{packageName}`) to avoid high-cardinality labels
- **Partial failure tolerance**: `FetchAll` continues with successful indexes if some fail
- **Atomic writes**: Publisher uses temp file + rename to prevent partial reads

## Build & Test

```bash
go build -o apk-datasource ./cmd/apk-datasource
go test ./...
go vet ./...
```

Container image built with Chainguard wolfi-base (multi-stage, nonroot user). Dockerfile at `deploy/docker/Dockerfile`.

## Version Injection

ldflags set version, commit, date at build time:
```
-X github.com/cwaits6/apk-datasource/cmd/apk-datasource/cmd.version=...
-X github.com/cwaits6/apk-datasource/cmd/apk-datasource/cmd.commit=...
-X github.com/cwaits6/apk-datasource/cmd/apk-datasource/cmd.date=...
```

## Commit Conventions

Uses [Conventional Commits](https://www.conventionalcommits.org/) with semantic-release:
- `feat:` — new feature (minor bump)
- `fix:` — bug fix (patch bump)
- `feat!:` or `BREAKING CHANGE:` — major bump
- `chore:` / `docs:` — no release

## CI/CD

- **ci.yml** — lint (golangci-lint), test, build on PRs
- **pr.yml** — dependency-review, trivy, semgrep (reuses cwaits6/.github shared workflows)
- **release.yml** — semantic-release on push to main (Go binary builds via shared workflow)
- **container-build.yml** — multi-arch (amd64/arm64) buildah builds; GHCR always, Docker Hub on tags only
- **generate-index.yml** — every 4h, generates static index and deploys to GitHub Pages
- **check-go-version.yml** — weekly check for newer Go minor version in Wolfi, auto-creates PR
- **scorecard.yml** — OpenSSF Scorecard weekly + on push to main

Shared workflows live in `cwaits6/.github` repo.

## Deploy

- **Docker Compose**: `deploy/docker/docker-compose.yml` — uses published GHCR image
- **Helm**: `charts/apk-datasource/` — deployment, service, optional ingress/cronjob, Prometheus annotations when `metrics.enabled`
- **GitLab CI**: `deploy/gitlab/.gitlab-ci.yml` — build + GitLab Pages deploy
- **GitHub Pages**: Hosted index at `https://cwaits6.github.io/apk-datasource/`

## Wolfi-Specific

- Wolfi uses dash for Go package names: `go-1.26` not `go~1.26`
- Dockerfile pins with `apk add go-1.26=<version>` — see `deploy/docker/Dockerfile` for the current exact pin
- APKINDEX URL pattern: `https://packages.wolfi.dev/os/{arch}/APKINDEX.tar.gz`
