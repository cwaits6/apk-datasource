# Renovate Configuration for apk-datasource

## Setup

Add the following to your `renovate.json` (or `renovate.json5`) to auto-update pinned APK package versions in Dockerfiles.

### Using the hosted index (GitHub Pages)

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
      "matchStrings": [
        "apk add[^\\n]*?\\s(?<depName>[a-zA-Z0-9][a-zA-Z0-9_.+\\-]*)=(?<currentValue>[^\\s]+)"
      ],
      "datasourceTemplate": "custom.apk-wolfi",
      "versioningTemplate": "loose"
    }
  ]
}
```

### Using a self-hosted instance

Replace the `defaultRegistryUrlTemplate` with your server URL:

```json
{
  "customDatasources": {
    "apk-wolfi": {
      "defaultRegistryUrlTemplate": "https://your-server.example.com/x86_64/{{packageName}}.json",
      "format": "json"
    }
  }
}
```

### Multi-architecture support

If you build for multiple architectures, create separate datasources:

```json
{
  "customDatasources": {
    "apk-wolfi-x86": {
      "defaultRegistryUrlTemplate": "https://cwaits6.github.io/apk-datasource/x86_64/{{packageName}}.json",
      "format": "json"
    },
    "apk-wolfi-arm": {
      "defaultRegistryUrlTemplate": "https://cwaits6.github.io/apk-datasource/aarch64/{{packageName}}.json",
      "format": "json"
    }
  }
}
```

## How it works

1. The `customManagers` regex matches `apk add package=version` patterns in Dockerfiles
2. Renovate queries the `customDatasources` URL with the package name
3. The JSON response contains all available versions
4. Renovate creates a PR if a newer version is available
