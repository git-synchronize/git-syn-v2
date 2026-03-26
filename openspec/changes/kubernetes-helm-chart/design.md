## Context

Kubernetes users need a production-ready way to host a git mirror. A Helm chart packages the required Kubernetes resources — Deployment, Service, Ingress, PVC — into a single installable artifact. cert-manager with Let's Encrypt is the de-facto standard for automatic TLS on Kubernetes, so the chart targets that stack.

## Goals / Non-Goals

**Goals:**
- Standard Helm chart structure installable with `helm install`
- Configurable ingress hostname, storage size, replica count, and image tag via `values.yaml`
- Let's Encrypt TLS via cert-manager annotation on the Ingress
- PVC for persistent repository storage

**Non-Goals:**
- Web UI or repository browser
- RBAC resources (ServiceAccount, ClusterRole)
- Horizontal pod autoscaling
- cert-manager installation (documented as a prerequisite)

## Decisions

**Single Deployment running the git-http-backend container.**
Mirrors the Docker Compose design. One pod serves all repositories from the mounted PVC. Replica count defaults to 1 (a git server with a single PVC can't safely scale to >1 without a ReadWriteMany storage class).

**Ingress uses `kubernetes.io/ingress.class: nginx` and `cert-manager.io/cluster-issuer` annotation.**
nginx-ingress + cert-manager is the most common production stack. Alternative: Traefik-specific annotations. Rejected — nginx is more universally available. Users on Traefik can override the annotations via `values.yaml`.

**PVC uses `ReadWriteOnce` access mode.**
Matches the single-replica constraint. `ReadWriteMany` would require a NFS-backed storage class which is not universally available.

**`values.yaml` exposes `ingressHost` as a required field (empty default, chart fails if not set).**
An ingress without a hostname is useless. Making it required surfaces the misconfiguration at install time rather than silently creating a broken Ingress resource.

## Risks / Trade-offs

[ReadWriteOnce limits availability] → A pod reschedule will cause brief downtime while the PVC is re-attached. Acceptable for a DR mirror that is not in the critical path of development. Document this limitation.

[cert-manager dependency] → Users without cert-manager must manually provision TLS or disable the Ingress TLS section. Provide a `tls.enabled` value to make TLS optional.
