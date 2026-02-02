package analysis

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/jeanlopezxyz/mcp-prometheus/pkg/mcputil"
	"github.com/jeanlopezxyz/mcp-prometheus/pkg/prometheus"
)

func registerInvestigatePod(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "investigatePod",
		Description: "Deep investigation of a specific pod: status, resources, restarts, OOM events, container details.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Analysis: Investigate Pod",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"podName": {
					Type:        "string",
					Description: "Pod name to investigate",
				},
				"namespace": {
					Type:        "string",
					Description: "Namespace of the pod",
				},
			},
			Required: []string{"podName", "namespace"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := mcputil.GetArguments(request)
		if err != nil {
			return mcputil.NewErrorResult(err.Error()), nil
		}
		podName, _ := args["podName"].(string)
		if podName == "" {
			return mcputil.NewErrorResult("podName parameter is required"), nil
		}
		namespace, _ := args["namespace"].(string)
		if namespace == "" {
			return mcputil.NewErrorResult("namespace parameter is required"), nil
		}

		queries := map[string]string{
			"Pod Phase":            fmt.Sprintf(`kube_pod_status_phase{namespace="%s",pod="%s"}`, namespace, podName),
			"Container Status":     fmt.Sprintf(`kube_pod_container_status_running{namespace="%s",pod="%s"}`, namespace, podName),
			"CPU Usage (cores)":    fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s",container!=""}[5m]))`, namespace, podName),
			"Memory Usage (bytes)": fmt.Sprintf(`sum(container_memory_working_set_bytes{namespace="%s",pod="%s",container!=""})`, namespace, podName),
			"CPU Requests":         fmt.Sprintf(`sum(kube_pod_container_resource_requests{namespace="%s",pod="%s",resource="cpu"})`, namespace, podName),
			"CPU Limits":           fmt.Sprintf(`sum(kube_pod_container_resource_limits{namespace="%s",pod="%s",resource="cpu"})`, namespace, podName),
			"Memory Requests":      fmt.Sprintf(`sum(kube_pod_container_resource_requests{namespace="%s",pod="%s",resource="memory"})`, namespace, podName),
			"Memory Limits":        fmt.Sprintf(`sum(kube_pod_container_resource_limits{namespace="%s",pod="%s",resource="memory"})`, namespace, podName),
			"Restarts":             fmt.Sprintf(`sum(kube_pod_container_status_restarts_total{namespace="%s",pod="%s"})`, namespace, podName),
			"OOM Killed":           fmt.Sprintf(`kube_pod_container_status_last_terminated_reason{namespace="%s",pod="%s",reason="OOMKilled"}`, namespace, podName),
			"Network RX (bytes/s)": fmt.Sprintf(`sum(rate(container_network_receive_bytes_total{namespace="%s",pod="%s"}[5m]))`, namespace, podName),
			"Network TX (bytes/s)": fmt.Sprintf(`sum(rate(container_network_transmit_bytes_total{namespace="%s",pod="%s"}[5m]))`, namespace, podName),
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("=== Pod Investigation: %s/%s ===\n\n", namespace, podName))

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
