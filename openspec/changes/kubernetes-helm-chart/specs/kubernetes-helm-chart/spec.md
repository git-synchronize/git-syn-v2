## ADDED Requirements

### Requirement: Helm chart structure
The `deploy/helm/git-syn/` directory SHALL contain a valid Helm chart installable with `helm install`.

#### Scenario: Chart installs successfully
- **WHEN** the user runs `helm install git-syn ./deploy/helm/git-syn --set ingressHost=git.example.com`
- **THEN** the chart SHALL create a Deployment, Service, Ingress, and PersistentVolumeClaim in the target namespace
- **THEN** `helm lint` SHALL pass without errors

### Requirement: Configurable values
The chart's `values.yaml` SHALL expose key deployment parameters.

#### Scenario: Default values applied
- **WHEN** the chart is installed without overrides
- **THEN** `replicaCount` SHALL default to `1`
- **THEN** `storage` SHALL default to `10Gi`
- **THEN** `image.tag` SHALL default to `latest`

#### Scenario: ingressHost is required
- **WHEN** the chart is installed without setting `ingressHost`
- **THEN** the Ingress resource SHALL be created with an empty host (behavior depends on ingress controller) and the README SHALL warn that `ingressHost` must be set

### Requirement: Let's Encrypt TLS via cert-manager
The Ingress resource SHALL include a cert-manager annotation for automatic TLS certificate provisioning.

#### Scenario: TLS enabled (default)
- **WHEN** `tls.enabled` is `true` (default)
- **THEN** the Ingress SHALL include the `cert-manager.io/cluster-issuer` annotation
- **THEN** the Ingress SHALL include a `tls` block referencing the generated certificate secret

#### Scenario: TLS disabled
- **WHEN** `tls.enabled` is set to `false`
- **THEN** the Ingress SHALL be created without TLS configuration

### Requirement: Persistent storage
The chart SHALL create a PVC for repository data.

#### Scenario: PVC created on install
- **WHEN** the chart is installed
- **THEN** a PersistentVolumeClaim SHALL be created with `ReadWriteOnce` access mode and the configured storage size
- **THEN** the Deployment SHALL mount the PVC at the repository data path
