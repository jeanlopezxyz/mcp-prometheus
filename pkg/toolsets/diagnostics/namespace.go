package diagnostics

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeanlopezxyz/mcp-prometheus/pkg/mcputil"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
)

func registerDiagnoseNamespace(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "diagnoseNamespace",
		Description: "Diagnose a namespace's health: pod status, resource consumption, restarts, storage, network.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Diagnostics: Namespace",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"namespace": {
					Type:        "string",
					Description: "Namespace to diagnose",
				},
			},
			Required: []string{"namespace"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := mcputil.GetArguments(request)
		if err != nil {
			return mcputil.NewErrorResult(err.Error()), nil
		}
		namespace, _ := args["namespace"].(string)
		if namespace == "" {
			return mcputil.NewErrorResult("namespace parameter is required"), nil
		}

		queries := map[string]string{
			"Running Pods":         fmt.Sprintf(`count(kube_pod_status_phase{namespace="%s",phase="Running"})`, namespace),
			"Pending Pods":         fmt.Sprintf(`count(kube_pod_status_phase{namespace="%s",phase="Pending"})`, namespace),
			"Failed Pods":          fmt.Sprintf(`count(kube_pod_status_phase{namespace="%s",phase="Failed"})`, namespace),
			"CPU Usage (cores)":    fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{namespace="%s",container!=""}[5m]))`, namespace),
			"Memory Usage (bytes)": fmt.Sprintf(`sum(container_memory_working_set_bytes{namespace="%s",container!=""})`, namespace),
			"Restarts (1h)":        fmt.Sprintf(`sum(increase(kube_pod_container_status_restarts_total{namespace="%s"}[1h]))`, namespace),
			"Network RX (bytes/s)": fmt.Sprintf(`sum(rate(container_network_receive_bytes_total{namespace="%s"}[5m]))`, namespace),
			"Network TX (bytes/s)": fmt.Sprintf(`sum(rate(container_network_transmit_bytes_total{namespace="%s"}[5m]))`, namespace),
			"PVC Usage":            fmt.Sprintf(`kubelet_volume_stats_used_bytes{namespace="%s"}`, namespace),
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("=== Namespace Diagnostics: %s ===\n\n", namespace))

		for label, query := range queries {
			result, err := client.Query(query)
			if err != nil {
				sb.WriteString(fmt.Sprintf("%s: error - %v\n", label, err))
			} else {
				sb.WriteString(fmt.Sprintf("%s:\n%s\n\n", label, result))
			}
		}

		return mcputil.NewTextResult(sb.String()), nil
	})
}
