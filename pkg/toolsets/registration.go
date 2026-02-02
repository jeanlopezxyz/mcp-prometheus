package toolsets

import (
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/toolsets/analysis"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/toolsets/basic"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/toolsets/diagnostics"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all Prometheus MCP tools with the server.
func RegisterAll(s *mcp.Server, client *prometheus.Client) {
	basic.Register(s, client)
	diagnostics.Register(s, client)
	analysis.Register(s, client)
}
