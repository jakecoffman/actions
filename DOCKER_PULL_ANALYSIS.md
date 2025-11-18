# Docker Pull HTTP Request Analysis Workflows

This repository includes two GitHub Actions workflows designed to capture and count HTTP requests made during a `docker pull` operation.

## Workflows

### 1. Docker Debug Pull (`docker-debug-pull.yml`)

This workflow uses Docker's debug mode and network packet capture to analyze requests.

**Features:**
- Enables Docker daemon debug logging
- Captures network traffic using `tcpdump`
- Analyzes both Docker logs and network packets
- Uploads comprehensive logs as artifacts

**Usage:**
1. Go to Actions tab in GitHub
2. Select "Docker Debug Pull" workflow
3. Click "Run workflow"
4. Wait for completion and download artifacts

**Note:** Since HTTPS traffic is encrypted, this method may have limited visibility into actual HTTP methods.

### 2. Docker MITM Proxy Pull (`docker-mitm-pull.yml`)

This workflow uses mitmproxy as a man-in-the-middle proxy to capture all HTTP/HTTPS requests.

**Features:**
- Installs and configures mitmproxy with CA certificates
- Routes all Docker traffic through the proxy
- Captures unencrypted HTTP method and path information
- Provides detailed counts of HEAD and GET requests
- Specifically counts requests to Docker Registry v2 API endpoints (`/v2/`)

**Usage:**
1. Go to Actions tab in GitHub
2. Select "Docker MITM Proxy Pull" workflow
3. Click "Run workflow"
4. Wait for completion
5. Review the analysis in the workflow logs
6. Download detailed logs from artifacts

**Output Example:**
```
SUMMARY OF DOCKER REGISTRY (/v2/) REQUESTS:
  HEAD /v2/ requests: X
  GET /v2/ requests: Y
  Total /v2/ requests: X+Y
```

## Target Image

Both workflows pull the same image:
```
ghcr.io/dependabot/dependabot-updater-npm:641ac3aeae97764683fb989a25dc3fa9340b7e24
```

## Recommendation

For the most accurate HTTP request counting, use the **Docker MITM Proxy Pull** workflow, as it can decrypt HTTPS traffic and provide exact method and path information.

## Artifacts

Both workflows upload their logs as artifacts that persist for 7 days:
- `docker-debug-analysis` - Contains Docker logs and packet captures
- `mitm-proxy-analysis` - Contains detailed mitmproxy request logs
