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

// Register registers all diagnostic tools.
func Register(s *mcp.Server, client *prometheus.Client) {
	registerClusterHealth(s, client)
	registerDiagnoseNode(s, client)
	registerDiagnoseNamespace(s, client)
}

func registerClusterHealth(s *mcp.Server, client *prometheus.Client) {
	s.AddTool(&mcp.Tool{
		Name:        "getClusterHealthOverview",
		Description: "Get comprehensive cluster health overview. Combines node status, pod health, resource usage, and active alerts. Use as first step for troubleshooting.",
		Annotations: &mcp.ToolAnnotations{
			Title:        "Diagnostics: Cluster Health",
			ReadOnlyHint: true,
		},
		InputSchema: &jsonschema.Schema{
			Type: "object",
		},
	}, func(ctx context.Context, request *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		queries := map[string]string{
			"Node Count":             `count(kube_node_info)`,
			"Ready Nodes":            `count(kube_node_status_condition{condition="Ready",status="true"})`,
			"Running Pods":           `count(kube_pod_status_phase{phase="Running"})`,
			"Pending Pods":           `count(kube_pod_status_phase{phase="Pending"})`,
			"Failed Pods":            `count(kube_pod_status_phase{phase="Failed"})`,
			"Cluster CPU Usage %":    `100 * (1 - avg(rate(node_cpu_seconds_total{mode="idle"}[5m])))`,
			"Cluster Memory Usage %": `100 * (1 - sum(node_memory_MemAvailable_bytes) / sum(node_memory_MemTotal_bytes))`,
			"Pod Restart Count (1h)": `sum(increase(kube_pod_container_status_restarts_total[1h]))`,
			"Firing Alerts":          `count(ALERTS{alertstate="firing"})`,
		}

		var sb strings.Builder
		sb.WriteString("=== Cluster Health Overview ===\n\n")

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
