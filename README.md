# Prometheus MCP Server

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![npm version](https://img.shields.io/npm/v/mcp-prometheus)](https://www.npmjs.com/package/mcp-prometheus)
[![Java](https://img.shields.io/badge/Java-21+-orange)](https://adoptium.net/)
[![GitHub release](https://img.shields.io/github/v/release/jeanlopezxyz/mcp-prometheus)](https://github.com/jeanlopezxyz/mcp-prometheus/releases/latest)

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for Prometheus integration.

Built with [Quarkus MCP Server](https://docs.quarkiverse.io/quarkus-mcp-server/dev/index.html).

## Transport Modes

| Mode | Description | Use Case |
|------|-------------|----------|
| **stdio** | Standard input/output | Default for Claude Code, Claude Desktop, Cursor, VS Code |
| **SSE** | Server-Sent Events over HTTP | Standalone server, multiple clients |

## Requirements

- **Java 21+** - [Download](https://adoptium.net/)
- **Prometheus** - Running and accessible

---

## Installation

### Quick Install (Claude Code CLI)

```bash
claude mcp add prometheus -e PROMETHEUS_URL="http://localhost:9090" -- npx -y mcp-prometheus@latest
```

### Claude Code

Add to `~/.claude/settings.json`:

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

### Claude Desktop

Add to `claude_desktop_config.json`:

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

### SSE Mode

```bash
PROMETHEUS_URL="http://localhost:9090" npx mcp-prometheus --port 9081
```

Endpoint: `http://localhost:9081/mcp/sse`

---

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PROMETHEUS_URL` | Prometheus API URL | `http://localhost:9090` |

### Command Line Options

| Option | Description |
|--------|-------------|
| `--port <PORT>` | Start in SSE mode on specified port |
| `--help` | Show help message |
| `--version` | Show version |

---

## Tools

This server provides **5 tools**:

### `query`
Execute a PromQL query. Returns current metric values.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `promql` | string | Yes | PromQL query to execute |

**Examples:**
- Check targets: `query promql='up'`
- CPU usage: `query promql='rate(node_cpu_seconds_total{mode="idle"}[5m])'`
- Memory: `query promql='node_memory_MemAvailable_bytes'`

---

### `queryRange`
Execute a range PromQL query. Returns metric values over time.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `promql` | string | Yes | PromQL query |
| `duration` | string | Yes | Time duration: `1h`, `30m`, `24h`, `7d` |
| `step` | string | No | Step interval: `1m`, `5m` (default: `1m`) |

**Example:**
- CPU over 1 hour: `queryRange promql='rate(node_cpu_seconds_total[5m])' duration='1h' step='5m'`

---

### `getTargets`
Get Prometheus scrape targets status.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `state` | string | No | Filter: `active`, `dropped`, `any` (default: `any`) |

---

### `getRules`
Get Prometheus alerting and recording rules.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | string | No | Type: `alerting`, `recording`, `all` (default: `all`) |

---

### `getPrometheusStatus`
Get Prometheus server status: version, build info, and runtime.

**Parameters:** None

---

## Example Prompts

Use natural language to query Prometheus. Here are prompts organized by use case:

### Quick Health Checks

```
"Are all services up?"
"Show me which targets are healthy"
"What services are down?"
"Check if the API is responding"
"Is Prometheus scraping all targets successfully?"
```

### CPU Monitoring

```
"What's the current CPU usage?"
"Show CPU usage across all nodes"
"Which pods are using the most CPU?"
"Show me CPU usage for the last hour"
"Is any container hitting CPU limits?"
"What's the CPU usage trend over the past 6 hours?"
```

### Memory Analysis

```
"How much memory is available on each node?"
"Show memory usage across the cluster"
"Which pods are using the most memory?"
"Are any containers close to their memory limits?"
"Show me memory usage trends for the database"
"What's the memory consumption in the production namespace?"
```

### Kubernetes Monitoring

```
"How many pods are running?"
"Show pods by namespace"
"Are any pods in CrashLoopBackOff?"
"What's the replica count for the web deployment?"
"Show me pending pods"
"How many nodes are in the cluster?"
"What's the pod distribution across nodes?"
```

### Request & Latency Metrics

```
"What's the request rate for the API?"
"Show HTTP error rates"
"What's the 99th percentile latency?"
"Show me request duration over the last hour"
"Are there any 5xx errors?"
"What's the traffic pattern for the last 24 hours?"
```

### Disk & Storage

```
"How much disk space is available?"
"Show disk usage across all nodes"
"Which persistent volumes are running low?"
"What's the disk I/O rate?"
"Show me storage trends for the database volume"
```

### Alerting Rules

```
"What alerting rules are defined?"
"Which alerts are currently firing in Prometheus?"
"Show me pending alerts"
"What are the thresholds for memory alerts?"
"List all recording rules"
```

### Historical Analysis

```
"Show me CPU usage for the past week"
"What was the memory consumption yesterday?"
"Graph request latency over the last 24 hours"
"When did the error rate spike?"
"Compare today's traffic to yesterday"
```

### Troubleshooting

```
"Why might the API be slow? Show me relevant metrics"
"Investigate high memory usage on node-1"
"Show me all metrics for the payment service"
"What changed in the last hour? Something broke"
"Help me understand the current resource usage"
```

### Custom PromQL Queries

```
"Run this query: rate(http_requests_total[5m])"
"Execute: sum by (namespace) (kube_pod_info)"
"Query the total number of requests in the last hour"
"Calculate the error percentage for the API"
"Show me the top 5 pods by memory usage"
```

---

## Development

### Run in dev mode

```bash
export PROMETHEUS_URL="http://localhost:9090"
./mvnw quarkus:dev
```

### Build

```bash
./mvnw package -DskipTests -Dquarkus.package.jar.type=uber-jar
```

---

## License

[MIT](LICENSE) - Free to use, modify, and distribute.


