# Self-Hosting apk-datasource

## GitHub Pages

The included GitHub Actions workflow (`.github/workflows/generate-index.yml`) generates the index every 4 hours and deploys to GitHub Pages automatically. Fork the repo and enable GitHub Pages on the `gh-pages` branch.

## GitLab Pages

Use the included `.gitlab-ci.yml` in `deploy/gitlab/`:

```bash
cp deploy/gitlab/.gitlab-ci.yml .gitlab-ci.yml
```

Push to your GitLab repo and the pipeline will generate the index and deploy to GitLab Pages.

## Docker Compose

```bash
docker compose -f deploy/docker/docker-compose.yml up -d
```

The server runs on port 3000 with a 4-hour refresh interval. Edit `docker-compose.yml` to customize index URLs or other settings.

## Helm Chart

```bash
helm install apk-datasource ./charts/apk-datasource \
  --set serve.indexURLs[0]=https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz
```

### With ingress

```bash
helm install apk-datasource ./charts/apk-datasource \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=apk.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

### CronJob mode

Instead of the live server, generate static files on a schedule:

```yaml
cronjob:
  enabled: true
  schedule: "0 */4 * * *"
  persistence:
    enabled: true
    size: 1Gi
```

## Standalone server

```bash
apk-index serve \
  --index-url https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz \
  --port 3000 \
  --refresh-interval 4h
```

## S3 (planned)

S3 publishing is not yet implemented. For now, use the `generate` command and sync the output directory to S3:

```bash
apk-index generate \
  --index-url https://packages.wolfi.dev/os/x86_64/APKINDEX.tar.gz \
  --output-dir ./output

aws s3 sync ./output s3://your-bucket/apk-datasource/ --delete
```
