# apk-datasource

Generate [Renovate-compatible custom datasource](https://docs.renovatebot.com/modules/datasource/custom/) files from Wolfi and Alpine APK package indexes.

This solves the problem of pinning and auto-updating APK package versions in Dockerfiles — see [Renovate issue #5422](https://github.com/renovatebot/renovate/issues/5422).

## How It Works

`apk-index` fetches `APKINDEX.tar.gz` archives from Wolfi or Alpine repositories, parses the package metadata, and outputs one JSON file per package conforming to Renovate's custom datasource schema:

```json
{
  "releases": [
    { "version": "8.11.1-r0" }
  ],
  "sourceUrl": "https://github.com/wolfi-dev/os",
  "homepage": "https://curl.se"
}
```

## Hosted Index

A pre-built index for Wolfi packages is published to GitHub Pages and refreshed every 4 hours:

```
https://cwaits6.github.io/apk-datasource/{arch}/{package}.json
```

Example: `https://cwaits6.github.io/apk-datasource/x86_64/curl.json`

See [examples/renovate.json](examples/renovate.json) for a ready-to-use Renovate config.

## Quickstart

### Generate static files

```bash
go install github.com/cwaits6/apk-datasource/cmd/apk-index@latest

apk-index generate \
  --index-url https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz \
  --output-dir ./output
```

### Serve over HTTP

```bash
apk-index serve \
  --index-url https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz \
  --index-url https://packages.wolfi.dev/os/aarch64/APKINDEX.tar.gz \
  --port 3000
```

### Docker Compose

```bash
docker compose -f deploy/docker/docker-compose.yml up
```

## Self-Hosting

See [docs/SELF-HOSTING.md](docs/SELF-HOSTING.md) for GitHub Pages, GitLab Pages, S3, Helm, and Docker Compose options.

## Renovate Configuration

See [docs/RENOVATE-CONFIG.md](docs/RENOVATE-CONFIG.md) for complete setup instructions.

## CLI Reference

```
apk-index generate  # Fetch indexes and write JSON files to disk
apk-index serve     # Serve JSON over HTTP with periodic refresh
apk-index version   # Print version info
```

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

## License

[Apache-2.0](LICENSE)
