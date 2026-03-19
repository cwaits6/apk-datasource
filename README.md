<h1 align="center">apk-datasource</h1>

<p align="center">
  <img src="https://img.shields.io/github/go-mod/go-version/cwaits6/apk-datasource?logo=go" alt="Go Version">
  <a href="https://github.com/cwaits6/apk-datasource/releases/latest"><img src="https://img.shields.io/github/v/release/cwaits6/apk-datasource?logo=github" alt="Release"></a>
  <a href="https://github.com/cwaits6/apk-datasource/actions/workflows/ci.yml"><img src="https://github.com/cwaits6/apk-datasource/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/cwaits6/apk-datasource/actions/workflows/container-build.yml"><img src="https://github.com/cwaits6/apk-datasource/actions/workflows/container-build.yml/badge.svg" alt="Container Build"></a>
  <a href="https://scorecard.dev/viewer/?uri=github.com/cwaits6/apk-datasource"><img src="https://api.scorecard.dev/projects/github.com/cwaits6/apk-datasource/badge" alt="OpenSSF Scorecard"></a>
  <a href="https://github.com/cwaits6/apk-datasource/pkgs/container/apk-datasource"><img src="https://img.shields.io/badge/ghcr.io-cwaits6%2Fapk--datasource-blue?logo=github" alt="GHCR"></a>
  <a href="https://hub.docker.com/r/cwaits6/apk-datasource"><img src="https://img.shields.io/docker/v/cwaits6/apk-datasource?logo=docker&label=docker%20hub" alt="Docker Hub"></a>
  <a href="https://github.com/cwaits6/apk-datasource/blob/main/LICENSE"><img src="https://img.shields.io/github/license/cwaits6/apk-datasource" alt="License"></a>
</p>

Auto-update pinned APK package versions in Dockerfiles using [Renovate](https://docs.renovatebot.com/). Supports Wolfi and Alpine package indexes.

Renovate has no built-in datasource for APK packages, so `apk add curl=8.11.1-r0` in your Dockerfiles can't be auto-updated. This project solves that. See [renovatebot/renovate#5422](https://github.com/renovatebot/renovate/issues/5422).

## Renovate Setup

Add this to your `renovate.json` and Renovate will start opening PRs to update pinned APK versions (e.g. `curl=8.15.0-r2` -> `curl=8.19.0-r1`):

```json
{ "extends": ["github>cwaits6/renovate-config"] }
```

That's it. This uses the [hosted index](#hosted-index) on GitHub Pages — no server to run.

<details>
<summary>Manual configuration (if you want full control)</summary>

```json
{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "customDatasources": {
    "apk-wolfi": {
      "defaultRegistryUrlTemplate": "https://cwaits6.github.io/apk-datasource/x86_64/{{packageName}}.json",
      "format": "json"
    }
  },
  "customManagers": [
    {
      "customType": "regex",
      "fileMatch": ["(^|/)Dockerfile[^/]*$"],
      "matchStringsStrategy": "recursive",
      "matchStrings": [
        "apk\\s+add[^\\n\\\\]*(?:\\\\[^\\S\\n]*\\n[^\\n\\\\]*)*",
        "(?<depName>[a-zA-Z0-9][a-zA-Z0-9._+-]*)=(?<currentValue>\\d+[^\\s\\\\]+)"
      ],
      "datasourceTemplate": "custom.apk-wolfi",
      "versioningTemplate": "loose"
    }
  ]
}
```

For a self-hosted instance, replace the `defaultRegistryUrlTemplate` URL with your server address (e.g. `http://apk-datasource:3000/x86_64/{{packageName}}`).

</details>

## Hosted Index

A pre-built index for Wolfi x86_64 packages is published to GitHub Pages and refreshed every 4 hours:

```text
https://cwaits6.github.io/apk-datasource/x86_64/{package}.json
```

```bash
curl -s https://cwaits6.github.io/apk-datasource/x86_64/curl.json | jq .
```

## Quick Start

### Install

```bash
go install github.com/cwaits6/apk-datasource/cmd/apk-datasource@latest
```

### Generate static files

```bash
apk-datasource generate \
  --index-url https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz \
  --output-dir ./output
```

### Serve over HTTP

```bash
apk-datasource serve \
  --index-url https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz \
  --port 3000
```

## Deployment

### Docker Compose

```bash
docker compose -f deploy/docker/docker-compose.yml up -d
```

The server runs on port 3000 with a 4-hour refresh interval. Edit `deploy/docker/docker-compose.yml` to customize index URLs or other settings.

### Helm

```bash
helm install apk-datasource ./charts/apk-datasource \
  --set serve.indexURLs[0]=https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz
```

See [`charts/apk-datasource/`](charts/apk-datasource/) for all configurable values.

## CLI Reference

| Command | Description |
|---------|-------------|
| `apk-datasource generate` | Fetch indexes and write JSON files to disk |
| `apk-datasource serve` | Serve JSON over HTTP with periodic refresh |
| `apk-datasource version` | Print version info |

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--log-level` | `info` | Log level (debug, info, warn, error) |
| `--log-format` | `text` | Log format (text, json) |

### `generate` Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--index-url` | *(required)* | APKINDEX.tar.gz URL (repeatable) |
| `--output-dir` | `./output` | Output directory |
| `--source-url` | *(auto-detect)* | Override source URL |
| `--homepage` | *(from index)* | Override homepage |

### `serve` Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--index-url` | *(required)* | APKINDEX.tar.gz URL (repeatable) |
| `--port` | `3000` | HTTP port |
| `--refresh-interval` | `4h` | Refresh interval |
| `--source-url` | *(auto-detect)* | Override source URL |
| `--homepage` | *(from index)* | Override homepage |
| `--metrics` | `true` | Enable Prometheus metrics on `/metrics` |

## Metrics

When `--metrics` is enabled (the default), the server exposes a Prometheus-compatible `/metrics` endpoint on the same port. Available metrics:

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `http_requests_total` | Counter | method, path, status_code | Total HTTP requests |
| `http_request_duration_seconds` | Histogram | method, path, status_code | Request latency |
| `refresh_total` | Counter | status | Index refresh attempts |
| `refresh_duration_seconds` | Histogram | status | Refresh latency |
| `refresh_packages` | Gauge | — | Package count after last refresh |
| `server_ready` | Gauge | — | Server readiness (0/1) |

Scrape with Prometheus:

```yaml
scrape_configs:
  - job_name: apk-datasource
    static_configs:
      - targets: ["localhost:3000"]
```

The Helm chart adds `prometheus.io/*` annotations automatically when `metrics.enabled` is `true`.

## How It Works

`apk-datasource` fetches `APKINDEX.tar.gz` archives from Wolfi or Alpine repositories, parses the package metadata, and outputs one JSON file per package conforming to Renovate's [custom datasource schema](https://docs.renovatebot.com/modules/datasource/custom/):

```json
{
  "releases": [
    { "version": "8.11.1-r0" }
  ],
  "sourceUrl": "https://github.com/wolfi-dev/os",
  "homepage": "https://curl.se"
}
```

Two modes: `generate` writes static JSON files to disk, `serve` runs an HTTP server with periodic refresh.

## Contributing

1. Fork the repo
2. Create a feature branch
3. Submit a pull request

## License

[Apache-2.0](LICENSE)
