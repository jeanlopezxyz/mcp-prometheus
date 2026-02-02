package analysis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeanlopezxyz/mcp-prometheus/pkg/mcputil"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
)

func registerCompareTimeRanges(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "compareTimeRanges",
		Description: "Compare metric values between two time periods. Useful for before/after analysis (e.g., after deployment).",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Analysis: Compare Time Ranges",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"promql": {
					Type:        "string",
					Description: "PromQL query expression",
				},
				"currentPeriod": {
					Type:        "string",
					Description: "Current period: 1h, 6h, 24h (default: 1h)",
				},
				"comparisonOffset": {
					Type:        "string",
					Description: "How far back to compare: 1h, 24h, 7d (default: 24h)",
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

		currentPeriod := "1h"
		if cp, ok := args["currentPeriod"].(string); ok && cp != "" {
			currentPeriod = cp
		}

		comparisonOffset := "24h"
		if co, ok := args["comparisonOffset"].(string); ok && co != "" {
			comparisonOffset = co
		}

		periodDur, err := parseDuration(currentPeriod)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Invalid currentPeriod: %v", err)), nil
		}
		offsetDur, err := parseDuration(comparisonOffset)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Invalid comparisonOffset: %v", err)), nil
		}

		now := time.Now()
		step := "1m"
		if periodDur > 6*time.Hour {
			step = "5m"
		}

		// Current period
		currentResult, err := client.QueryRange(promql, now.Add(-periodDur), now, step)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Current period query failed: %v", err)), nil
		}

		// Comparison period (same duration, shifted back by offset)
		compEnd := now.Add(-offsetDur)
		compStart := compEnd.Add(-periodDur)
		compResult, err := client.QueryRange(promql, compStart, compEnd, step)
		if err != nil {
			return mcputil.NewErrorResult(fmt.Sprintf("Comparison period query failed: %v", err)), nil
		}

		result := fmt.Sprintf(
			"=== Time Range Comparison ===\nQuery: %s\nCurrent: last %s\nComparison: %s ago (same duration)\n\n"+
				"--- Current Period ---\n%s\n\n--- Comparison Period ---\n%s",
			promql, currentPeriod, comparisonOffset, currentResult, compResult,
		)
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
