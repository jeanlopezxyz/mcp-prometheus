# Prometheus MCP Server (Go)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![npm version](https://img.shields.io/npm/v/mcp-prometheus)](https://www.npmjs.com/package/mcp-prometheus)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![GitHub release](https://img.shields.io/github/v/release/jeanlopezxyz/mcp-prometheus?sort=semver)](https://github.com/jeanlopezxyz/mcp-prometheus/releases/latest)

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for Prometheus integration. Native Go binary with built-in Kubernetes connectivity via client-go.

## Installation

### npx

```bash
npx -y mcp-prometheus@latest
```

### MCP Client Configuration

Add to your MCP client configuration (VS Code, Cursor, Windsurf, etc.):

```json
{
  "mcpServers": {
    "prometheus": {
      "command": "npx",
      "args": ["-y", "mcp-prometheus@latest"],
      "env": {
        "PROMETHEUS_URL": "http://localhost:9090"
      }
    }
  }
}
```

### VS Code

```shell
code --add-mcp '{"name":"prometheus","command":"npx","args":["-y","mcp-prometheus@latest"],"env":{"PROMETHEUS_URL":"http://localhost:9090"}}'
```

### Kubernetes Auto-Connect

Connects directly to Prometheus running in Kubernetes via the K8S API service proxy. Uses native kubeconfig/in-cluster config via client-go. No `kubectl` or port-forwarding required.

```json
{
  "mcpServers": {
    "prometheus": {
      "command": "npx",
      "args": ["-y", "mcp-prometheus@latest"],
      "env": {
        "K8S_NAMESPACE": "openshift-monitoring",
        "K8S_SERVICE": "prometheus-k8s"
      }
    }
  }
}
```

### Binary

Download from [GitHub Releases](https://github.com/jeanlopezxyz/mcp-prometheus/releases) or build from source:

```bash
make build
./mcp-prometheus
```

---

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PROMETHEUS_URL` | Prometheus API URL | `http://localhost:9090` |

### Kubernetes Auto-Connect Variables

| Variable | Required | Description | Default |
|----------|----------|-------------|---------|
| `K8S_NAMESPACE` | Yes | Kubernetes namespace of the service | - |
| `K8S_SERVICE` | Yes | Kubernetes service name | - |
| `K8S_SERVICE_PORT` | No | Service port | `9090` |
| `KUBECONFIG` | No | Path to kubeconfig file | `~/.kube/config` |

**Precedence:** `PROMETHEUS_URL` > `K8S_NAMESPACE` + `K8S_SERVICE` > `http://localhost:9090`

**Connection strategy:**
1. In-cluster config (when running inside a pod)
2. Kubeconfig file (`KUBECONFIG` or `~/.kube/config`)

---

## Tools (11)

### Basic Queries

| Tool | Description |
|------|-------------|
| `query` | Execute a PromQL instant query |
| `queryRange` | Execute a PromQL range query over time |
| `getTargets` | Get scrape targets status |
| `getRules` | Get alerting and recording rules |
| `getPrometheusStatus` | Get server version and runtime info |

### Cluster Diagnostics

| Tool | Description |
|------|-------------|
| `getClusterHealthOverview` | Comprehensive cluster health overview |
| `diagnoseNode` | Diagnose a specific node's health |
| `diagnoseNamespace` | Diagnose a namespace's health |

### Resource Analysis

| Tool | Description |
|------|-------------|
| `getTopResourceConsumers` | Top CPU/memory/network consumers |
| `investigatePod` | Deep investigation of a specific pod |
| `compareTimeRanges` | Compare metrics between two time periods |

---

## Example Prompts

```
"What's the current CPU usage across all nodes?"
"Show me the cluster health overview"
"Which pods are using the most memory?"
"Diagnose node worker-1"
"Are all scrape targets healthy?"
"Compare CPU usage now vs 24 hours ago"
"Investigate pod my-app in namespace production"
"What alerting rules are defined?"
```

---

## Development

### Build

```bash
make build              # Build for current platform
make build-all-platforms # Cross-compile for all platforms
```

### Container

```bash
podman build -f Containerfile -t mcp-prometheus .
```

---

## License

[MIT](LICENSE)
