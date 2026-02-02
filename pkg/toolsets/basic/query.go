package basic

import (
	"context"
	"fmt"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeanlopezxyz/mcp-prometheus/pkg/mcputil"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
)

func registerQuery(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "query",
		Description: "Execute a PromQL instant query against Prometheus. Returns current metric values.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Prometheus: Query",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"promql": {
					Type:        "string",
					Description: "PromQL query to execute",
				},
			},
			Required: []string{"promql"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := mcputil.GetArguments(request)
		if err != nil {
			return mcputil.NewErrorResult(err.Error()), nil
		}
		promql, _ := args["promql"].(string)
		if promql == "" {
			return mcputil.NewErrorResult("promql parameter is required"), nil
		}

		result, err := client.Query(promql)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Query failed: %v", err)), nil
		}
		return mcputil.NewTextResult(result), nil
	})
}

func registerQueryRange(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "queryRange",
		Description: "Execute a PromQL range query. Returns metric values over a time period.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Prometheus: Query Range",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"promql": {
					Type:        "string",
					Description: "PromQL query to execute",
				},
				"duration": {
					Type:        "string",
					Description: "Time range: 1h, 6h, 24h, 7d (default: 1h)",
				},
				"step": {
					Type:        "string",
					Description: "Step interval: 1m, 5m, 15m (default: 1m)",
				},
			},
			Required: []string{"promql"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := mcputil.GetArguments(request)
		if err != nil {
			return mcputil.NewErrorResult(err.Error()), nil
		}
		promql, _ := args["promql"].(string)
		if promql == "" {
			return mcputil.NewErrorResult("promql parameter is required"), nil
		}

		duration := "1h"
		if d, ok := args["duration"].(string); ok && d != "" {
			duration = d
		}
		step := "1m"
		if s, ok := args["step"].(string); ok && s != "" {
			step = s
		}

		dur, err := parseDuration(duration)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Invalid duration: %v", err)), nil
		}

		end := time.Now()
		start := end.Add(-dur)

		result, err := client.QueryRange(promql, start, end, step)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Range query failed: %v", err)), nil
		}
		return mcputil.NewTextResult(result), nil
	})
}

func parseDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration: %s", s)
	}
	unit := s[len(s)-1]
	valStr := s[:len(s)-1]
	var val int
	if _, err := fmt.Sscanf(valStr, "%d", &val); err != nil {
		return 0, fmt.Errorf("invalid duration value: %s", s)
	}
	switch unit {
	case 'm':
		return time.Duration(val) * time.Minute, nil
	case 'h':
		return time.Duration(val) * time.Hour, nil
	case 'd':
		return time.Duration(val) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unknown duration unit: %c", unit)
	}
}
