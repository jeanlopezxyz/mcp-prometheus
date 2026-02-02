package basic

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeanlopezxyz/mcp-prometheus/pkg/mcputil"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
)

// Register registers all basic Prometheus tools.
func Register(s *mcp.Server, client *prometheus.Client) {
	registerQuery(s, client)
	registerQueryRange(s, client)
	registerTargets(s, client)
	registerRules(s, client)
	registerStatus(s, client)
}

func registerStatus(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "getPrometheusStatus",
		Description: "Get Prometheus server status: version, build info, and runtime information.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Prometheus: Status",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		buildInfo, err := client.GetBuildInfo()
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Failed to get build info: %v", err)), nil
		}

		runtimeInfo, err := client.GetRuntimeInfo()
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Failed to get runtime info: %v", err)), nil
		}

		result := fmt.Sprintf("=== Build Info ===\n%s\n\n=== Runtime Info ===\n%s", buildInfo, runtimeInfo)
		return mcputil.NewTextResult(result), nil
	})
}
