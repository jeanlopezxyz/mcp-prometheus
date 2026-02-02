package analysis

import (
	"context"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"k8s.io/utils/ptr"

	"github.com/jeanlopezxyz/mcp-prometheus/pkg/mcputil"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
)

// Register registers all analysis tools.
func Register(s *mcp.Server, client *prometheus.Client) {
	registerTopResourceConsumers(s, client)
	registerInvestigatePod(s, client)
	registerCompareTimeRanges(s, client)
}

func registerTopResourceConsumers(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "getTopResourceConsumers",
		Description: "Get top resource consumers in the cluster. Identifies pods using most CPU, memory, or network.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Analysis: Top Resource Consumers",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"resourceType": {
					Type:        "string",
					Description: "Resource type: cpu, memory, network (default: cpu)",
				},
				"limit": {
					Type:        "number",
					Description: "Number of results 1-50 (default: 10)",
					Minimum:     ptr.To(float64(1)),
				},
			},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := mcputil.GetArguments(request)
		if err != nil {
			return mcputil.NewErrorResult(err.Error()), nil
		}
		resourceType := "cpu"
		if rt, ok := args["resourceType"].(string); ok && rt != "" {
			resourceType = rt
		}

		limit := 10
		if l, ok := args["limit"].(float64); ok && l > 0 {
			limit = int(l)
			if limit > 50 {
				limit = 50
			}
		}

		var query string
		switch resourceType {
		case "cpu":
			query = fmt.Sprintf(`topk(%d, sum by (namespace, pod) (rate(container_cpu_usage_seconds_total{container!=""}[5m])))`, limit)
		case "memory":
			query = fmt.Sprintf(`topk(%d, sum by (namespace, pod) (container_memory_working_set_bytes{container!=""}))`, limit)
		case "network":
			query = fmt.Sprintf(`topk(%d, sum by (namespace, pod) (rate(container_network_receive_bytes_total[5m]) + rate(container_network_transmit_bytes_total[5m])))`, limit)
		default:
			return mcputil.NewErrorResult(fmt.Sprintf("Unknown resource type: %s. Use cpu, memory, or network.", resourceType)), nil
		}

		result, err := client.Query(query)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Failed to get top consumers: %v", err)), nil
		}

		header := fmt.Sprintf("=== Top %d %s consumers ===\n\n", limit, resourceType)
		return mcputil.NewTextResult(header + result), nil
	})
}
