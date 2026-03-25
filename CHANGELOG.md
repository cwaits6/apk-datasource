# [1.1.0](https://github.com/cwaits6/apk-datasource/compare/v1.0.1...v1.1.0) (2026-03-25)


### Features

* default to Chainguard Wolfi index and make --index-url optional ([#12](https://github.com/cwaits6/apk-datasource/issues/12)) ([1ff4ee4](https://github.com/cwaits6/apk-datasource/commit/1ff4ee4a372b873b9ec2cbed76ebf106e431b94a))

## [1.0.1](https://github.com/cwaits6/apk-datasource/compare/v1.0.0...v1.0.1) (2026-03-19)


### Bug Fixes

* trigger container build on release event instead of tag push ([#7](https://github.com/cwaits6/apk-datasource/issues/7)) ([00a9b7a](https://github.com/cwaits6/apk-datasource/commit/00a9b7a3c23ec3a0d81e559e734d3d775b184eb3))

# [1.0.0](https://github.com/cwaits6/apk-datasource/compare/v0.1.0...v1.0.0) (2026-03-19)


* feat!: rename binary, add metrics, Docker Hub publishing, and README rewrite ([#6](https://github.com/cwaits6/apk-datasource/issues/6)) ([01263d5](https://github.com/cwaits6/apk-datasource/commit/01263d5bfb7ee1820aed25c8b904fcdd1112fc2f))


### Bug Fixes

* remove [skip ci] from semantic-release commit message ([#5](https://github.com/cwaits6/apk-datasource/issues/5)) ([50dd884](https://github.com/cwaits6/apk-datasource/commit/50dd8848f8b620bc583c38d892de3b6609f9363b))


### BREAKING CHANGES

* The binary is now named `apk-datasource` instead of
`apk-index`. Update any scripts, aliases, or CI pipelines that reference
the old name.

- Rename cmd/apk-index/ to cmd/apk-datasource/
- Update cobra Use field, version output, import paths
- Update Dockerfile build output, copy, entrypoint, and ldflags
- Update all GitHub and GitLab CI workflows
- Update .gitignore

* chore: clean up Docker Compose and Helm chart defaults

- Replace build block with published image in docker-compose.yml
- Add commented-out build block for local development
- Remove aarch64 index URL from Helm chart default values

* feat: add OpenTelemetry metrics with Prometheus /metrics endpoint

- Add pkg/metrics with OTel-Prometheus bridge (Setup/Noop pattern)
- Instruments: http_requests_total, http_request_duration_seconds,
  refresh_total, refresh_duration_seconds, refresh_packages, server_ready
- Add responseWriter wrapper and metricsMiddleware to server
- Add --metrics flag to serve command (default: true)
- Add Prometheus scrape annotations and --metrics arg to Helm chart
- Add metrics values section to Helm chart

* feat: publish container images to Docker Hub on tagged releases

Docker Hub receives only tagged versions (v*), not PR builds.
GHCR continues to receive both. Requires DOCKERHUB_USERNAME and
DOCKERHUB_TOKEN repository secrets and shared workflow support
for the dockerhub-image input.

* docs: rewrite README with badges and inline all documentation

- Add Go version, release, CI, OpenSSF Scorecard, and license badges
- Inline Renovate config and deployment docs from removed files
- Add metrics documentation and CLI reference for --metrics flag
- Remove docs/SELF-HOSTING.md, docs/RENOVATE-CONFIG.md, and
  examples/renovate.json (content now lives in README)

* feat: add OpenSSF Scorecard workflow

Runs weekly and on push to main. Publishes SARIF to GitHub Security
tab and results to scorecard.dev for badge support.

* docs: add logos to badges and container registry badges

* fix: address PR review feedback

- Normalize metric path labels to route templates to avoid high
  cardinality from unique package names
- Add nil receiver guards to Metrics methods
- Add test for /metrics endpoint
- Fix markdown lint: add language identifier to fenced code block
- Fix docker-compose.yml path reference in README

* fix: tune histogram buckets and remove noisy OTel scope labels

- Set HTTP duration buckets: 1ms to 10s (useful for percentiles)
- Set refresh duration buckets: 100ms to 60s
- Disable OTel scope info to remove empty otel_scope_schema_url and
  otel_scope_version labels from Prometheus output

* chore: add CLAUDE.md for AI-assisted development context

* fix: use tagged switch in normalizePath to satisfy staticcheck

* fix: address PR review nitpicks

- Remove hard-coded Go version pin from CLAUDE.md, point to Dockerfile
- Fix import ordering in server_test.go
- Wire server routes in metrics test to remove unused construction

# [0.1.0](https://github.com/cwaits6/apk-datasource/compare/v0.0.1...v0.1.0) (2026-03-18)


### Features

* build container images on pull requests ([#4](https://github.com/cwaits6/apk-datasource/issues/4)) ([ebbbfe5](https://github.com/cwaits6/apk-datasource/commit/ebbbfe50f5a41bebe1377dc27cd908374fed65a6))
* pass GitHub App secrets to release workflow ([#3](https://github.com/cwaits6/apk-datasource/issues/3)) ([1125bbd](https://github.com/cwaits6/apk-datasource/commit/1125bbd60c75b1e2c8e9ad5944b3428e20de08ea))

## [0.0.1](https://github.com/cwaits6/apk-datasource/compare/v0.0.0...v0.0.1) (2026-03-18)


### Bug Fixes

* add required permissions for release workflow ([#2](https://github.com/cwaits6/apk-datasource/issues/2)) ([175bcf5](https://github.com/cwaits6/apk-datasource/commit/175bcf5b38945238068bb1df268b8e3710ff16e3))
* fix workflow permissions and dockerfile apk package name ([#1](https://github.com/cwaits6/apk-datasource/issues/1)) ([9c84dca](https://github.com/cwaits6/apk-datasource/commit/9c84dca57502a79da83b482821edb0b326deb156))
