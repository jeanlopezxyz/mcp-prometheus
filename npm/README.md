# mcp-prometheus

MCP Server for Prometheus monitoring. Query metrics, analyze trends, and monitor infrastructure health through AI assistants.

## Quick Start

```bash
PROMETHEUS_URL="http://localhost:9090" npx mcp-prometheus
```

## Configuration

Add to `~/.claude/settings.json` (Claude Code) or your MCP client config:

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

## Requirements

- Java 21+
- Prometheus running and accessible

## Tools (5)

| Tool | Description |
|------|-------------|
| `query` | Execute instant PromQL query for current values |
| `queryRange` | Query metrics over time for trends/graphs |
| `getTargets` | Check scrape targets health and errors |
| `getRules` | List alerting and recording rules |
| `getPrometheusStatus` | Get server version and status |

## Example Prompts

```
"Are all services up?"
"What's the current CPU usage?"
"Show memory usage for the last hour"
"How many pods are running?"
"What alerting rules are defined?"
"Show me the request rate for the API"
"Which containers are using the most resources?"
"Run query: rate(http_requests_total[5m])"
```

## Documentation

Full docs: https://github.com/jeanlopezxyz/mcp-prometheus

## License

MIT
