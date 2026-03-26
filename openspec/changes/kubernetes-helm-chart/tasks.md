## 1. Chart scaffolding

- [x] 1.1 Create `deploy/helm/git-syn/` directory structure: `Chart.yaml`, `values.yaml`, `templates/`
- [x] 1.2 Write `Chart.yaml` with name, version (`0.1.0`), and description

## 2. values.yaml

- [x] 2.1 Define `replicaCount: 1`, `storage: 10Gi`, `ingressHost: ""`, `image.tag: latest`, `tls.enabled: true`, `tls.clusterIssuer: letsencrypt-prod`

## 3. Kubernetes templates

- [x] 3.1 Write `templates/deployment.yaml` for the git-serve Deployment mounting the PVC at the repository data path
- [x] 3.2 Write `templates/service.yaml` as a ClusterIP Service on port 80
- [x] 3.3 Write `templates/pvc.yaml` as a `ReadWriteOnce` PVC using the configured storage size
- [x] 3.4 Write `templates/ingress.yaml` with `cert-manager.io/cluster-issuer` annotation and TLS block (conditional on `tls.enabled`)

## 4. Documentation

- [x] 4.1 Write `deploy/helm/git-syn/README.md` documenting prerequisites (cert-manager), required values (`ingressHost`), and install command
- [x] 4.2 Document the `tls.enabled: false` escape hatch for clusters without cert-manager

## 5. Lint and dry-run

- [x] 5.1 Run `helm lint deploy/helm/git-syn` and resolve any warnings or errors
- [x] 5.2 Run `helm template git-syn deploy/helm/git-syn --set ingressHost=git.example.com` and verify all resources render correctly
