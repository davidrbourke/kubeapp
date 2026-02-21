# kubeapp

A CLI tool for inspecting Kubernetes clusters. Works with local clusters (minikube) and remote clusters via kubeconfig.

## Requirements

- Go 1.22+
- A kubeconfig file (default: `~/.kube/config`)

## Installation

```bash
git clone <repo>
cd kubeapp
GOTOOLCHAIN=local go build -o kubeapp ./cmd/cli/
```

The `kubeapp` binary will be created in the current directory.

## Usage

```
kubeapp [--kubeconfig <path>] [--namespace <ns>] <command>
```

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--kubeconfig` | `~/.kube/config` | Path to kubeconfig file |
| `-n, --namespace` | `default` | Target namespace |

---

### Pods

```bash
# List pods in the default namespace
kubeapp pods list

# List pods in a specific namespace
kubeapp pods list -n kube-system

# List pods across all namespaces
kubeapp pods list -A

# Describe a pod
kubeapp pods describe <pod-name>
kubeapp pods describe <pod-name> -n kube-system
```

**`pods list` output:**

```
NAME                                                 NAMESPACE       STATUS       READY
----                                                 ---------       ------       -----
coredns-76f75df574-abcde                             kube-system     Running      1/1
```

**`pods describe` output:**

```
Name:       coredns-76f75df574-abcde
Namespace:  kube-system
Node:       minikube
Status:     Running
Pod IP:     10.244.0.3
Started:    2024-01-01 00:00:00 +0000 UTC
Containers:
  Name:   coredns
  Image:  registry.k8s.io/coredns/coredns:v1.11.1
  Ports:  53/UDP, 53/TCP, 9153/TCP
  Ready:    true
  Restarts: 0
Conditions:
  PodScheduled         True
  ContainersReady      True
  Ready                True
```

---

### Services

```bash
# List services in the default namespace
kubeapp services list
kubeapp svc list          # alias

# List services in a specific namespace
kubeapp services list -n kube-system

# List services across all namespaces
kubeapp services list -A

# Describe a service
kubeapp services describe <service-name>
kubeapp services describe <service-name> -n kube-system
```

**`services list` output:**

```
NAME                                     NAMESPACE       TYPE         CLUSTER-IP       EXTERNAL-IP
----                                     ---------       ----         ----------       -----------
kubernetes                               default         ClusterIP    10.96.0.1        <none>
```

**`services describe` output:**

```
Name:       kubernetes
Namespace:  default
Type:       ClusterIP
ClusterIP:  10.96.0.1
Ports:
  <unnamed>         443:6443/TCP
Selector:
  component=apiserver
  provider=kubernetes
```

---

## Project Structure

```
cmd/cli/          CLI entry point and command definitions
  main.go         Root cobra command; k8s client initialisation
  pods.go         pods list / pods describe commands
  services.go     services list / services describe commands

internal/
  k8s/client.go         Kubernetes clientset factory
  pods/service.go       Pod list/get operations
  services/service.go   Service list/get operations
```

The `internal/` packages are decoupled from the CLI presentation layer, making them reusable by a future HTTP API (`cmd/api/`).

## Dependencies

- [client-go](https://github.com/kubernetes/client-go) v0.31.3
- [cobra](https://github.com/spf13/cobra) v1.8.1
