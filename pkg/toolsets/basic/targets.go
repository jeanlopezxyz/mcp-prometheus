package basic

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeanlopezxyz/mcp-prometheus/pkg/mcputil"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
)

func registerTargets(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "getTargets",
		Description: "Get Prometheus scrape targets status. Shows health, errors, last scrape time.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Prometheus: Targets",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"state": {
					Type:        "string",
					Description: "Filter by state: active, dropped",
				},
			},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := mcputil.GetArguments(request)
		if err != nil {
			return mcputil.NewErrorResult(err.Error()), nil
		}
		state, _ := args["state"].(string)

		result, err := client.GetTargets(state)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Failed to get targets: %v", err)), nil
		}
		return mcputil.NewTextResult(result), nil
	})
}
