package basic

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeanlopezxyz/mcp-prometheus/pkg/mcputil"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
)

func registerRules(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "getRules",
		Description: "Get Prometheus alerting and recording rules.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Prometheus: Rules",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"type": {
					Type:        "string",
					Description: "Filter by type: alert, record",
				},
			},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := mcputil.GetArguments(request)
		if err != nil {
			return mcputil.NewErrorResult(err.Error()), nil
		}
		ruleType, _ := args["type"].(string)

		result, err := client.GetRules(ruleType)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Failed to get rules: %v", err)), nil
		}
		return mcputil.NewTextResult(result), nil
	})
}
