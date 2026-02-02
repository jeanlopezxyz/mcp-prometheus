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

func registerDiagnoseNode(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "diagnoseNode",
		Description: "Diagnose a specific node's health: CPU, memory, disk, network, pods running, conditions.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Diagnostics: Node",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"nodeName": {
					Type:        "string",
					Description: "Node name to diagnose",
				},
			},
			Required: []string{"nodeName"},
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := mcputil.GetArguments(request)
		if err != nil {
			return mcputil.NewErrorResult(err.Error()), nil
		}
		nodeName, _ := args["nodeName"].(string)
		if nodeName == "" {
			return mcputil.NewErrorResult("nodeName parameter is required"), nil
		}

		queries := map[string]string{
			"CPU Usage %":          fmt.Sprintf(`100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle",instance=~"%s.*"}[5m])))`, nodeName),
			"Memory Usage %":       fmt.Sprintf(`100 * (1 - node_memory_MemAvailable_bytes{instance=~"%s.*"} / node_memory_MemTotal_bytes{instance=~"%s.*"})`, nodeName, nodeName),
			"Disk Usage %":         fmt.Sprintf(`100 * (1 - node_filesystem_avail_bytes{instance=~"%s.*",mountpoint="/"} / node_filesystem_size_bytes{instance=~"%s.*",mountpoint="/"})`, nodeName, nodeName),
			"Network RX (bytes/s)": fmt.Sprintf(`sum(rate(node_network_receive_bytes_total{instance=~"%s.*",device!="lo"}[5m]))`, nodeName),
			"Network TX (bytes/s)": fmt.Sprintf(`sum(rate(node_network_transmit_bytes_total{instance=~"%s.*",device!="lo"}[5m]))`, nodeName),
			"Pods Running":         fmt.Sprintf(`count(kube_pod_info{node="%s"})`, nodeName),
			"Node Conditions":      fmt.Sprintf(`kube_node_status_condition{node="%s",status="true"}`, nodeName),
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("=== Node Diagnostics: %s ===\n\n", nodeName))

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
