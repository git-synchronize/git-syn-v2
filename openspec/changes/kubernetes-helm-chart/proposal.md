## Why

Users running Kubernetes need a production-ready way to host a git mirror as a git-syn sync target. A Helm chart packages all required Kubernetes resources — storage, networking, TLS — into a single installable artifact.

## What Changes

- Add a Helm chart at `deploy/helm/git-syn/` with standard chart structure:
  - `Chart.yaml`: chart metadata (name, version, description)
  - `values.yaml`: configurable defaults (ingress hostname, storage size, replica count, image tag)
  - `templates/deployment.yaml`: `git-serve` deployment running git-http-backend
  - `templates/service.yaml`: ClusterIP service
  - `templates/ingress.yaml`: Ingress with cert-manager `cluster-issuer` annotation for Let's Encrypt TLS
  - `templates/pvc.yaml`: PersistentVolumeClaim for repository storage
- `values.yaml` defaults: `replicaCount: 1`, `storage: 10Gi`, `ingressHost: ""` (required), `image.tag: latest`

Out of scope: web UI, monitoring, RBAC beyond what Helm installs.

## Capabilities

### New Capabilities

- `kubernetes-helm-chart`: A Helm chart at `deploy/helm/git-syn/` for deploying a self-hosted git mirror on Kubernetes with Let's Encrypt TLS via cert-manager

### Modified Capabilities

<!-- none -->

## Impact

- `deploy/helm/git-syn/`: new directory with full Helm chart structure
- No Go code changes
- No new Go dependencies
- Requires cert-manager installed in target cluster (documented in chart README)
