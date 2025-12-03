package com.monitoring.prometheus.mcp;

import com.monitoring.prometheus.application.service.PrometheusService;
import io.quarkiverse.mcp.server.Tool;
import io.quarkiverse.mcp.server.ToolArg;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

/**
 * MCP Tools for Prometheus monitoring and observability.
 *
 * This server enables AI assistants to query metrics, analyze performance,
 * and monitor infrastructure health. Use it to check resource usage,
 * investigate issues, analyze trends, and understand alerting rules.
 *
 * Provides 5 tools in 4 categories:
 * - Queries: query, queryRange
 * - Targets: getTargets
 * - Rules: getRules
 * - Status: getPrometheusStatus
 */
@ApplicationScoped
public class PrometheusTools {

    @Inject
    PrometheusService prometheusService;

    // =========================================================================
    // Query Tools
    // =========================================================================

    @Tool(description = "Execute a PromQL instant query to get current metric values. "
            + "Use this to: check current CPU/memory usage, count running pods, verify service health, "
            + "get latest values for any metric. "
            + "Common queries: 'up' (target health), 'node_memory_MemAvailable_bytes' (free memory), "
            + "'rate(http_requests_total[5m])' (request rate), 'kube_pod_status_phase{phase=\"Running\"}' (running pods). "
            + "Returns point-in-time values.")
    public String query(
        @ToolArg(description = "PromQL query expression. Examples: 'up', 'node_cpu_seconds_total', "
                + "'rate(http_requests_total{job=\"api\"}[5m])', 'sum by (namespace) (kube_pod_info)'") String promql
    ) {
        return prometheusService.query(promql);
    }

    @Tool(description = "Execute a PromQL range query to get metric values over time. "
            + "Use this to: analyze trends, see historical data, investigate when issues started, "
            + "compare performance over time, prepare data for graphing. "
            + "Returns a time series of values at regular intervals over the specified duration.")
    public String queryRange(
        @ToolArg(description = "PromQL query expression. Same syntax as instant query.") String promql,
        @ToolArg(description = "How far back to query: '15m' (15 min), '1h' (1 hour), '6h', '24h', '7d' (7 days)") String duration,
        @ToolArg(description = "Resolution/granularity: '1m' (1 min), '5m', '15m', '1h'. Smaller = more data points. Default: 1m") String step
    ) {
        return prometheusService.queryRange(promql, duration, step);
    }

    // =========================================================================
    // Target Tools
    // =========================================================================

    @Tool(description = "Get Prometheus scrape target status. Shows all endpoints Prometheus is collecting metrics from, "
            + "their health status (up/down), last scrape time, and any errors. "
            + "Use this to: verify services are being monitored, diagnose collection issues, "
            + "find targets with scrape errors, check service discovery results.")
    public String getTargets(
        @ToolArg(description = "Filter targets: 'active' (currently being scraped), 'dropped' (removed by relabeling), 'any' (all). Default: any") String state
    ) {
        return prometheusService.getTargets(state);
    }

    // =========================================================================
    // Rules Tools
    // =========================================================================

    @Tool(description = "Get Prometheus alerting and recording rules. "
            + "Alerting rules define conditions that trigger alerts (sent to Alertmanager). "
            + "Recording rules pre-compute expensive queries for faster dashboards. "
            + "Use this to: see what alerts are defined, check rule evaluation status (firing/pending/inactive), "
            + "understand alert thresholds, review recording rule expressions.")
    public String getRules(
        @ToolArg(description = "Rule type: 'alerting' (alert definitions), 'recording' (pre-computed queries), 'all' (both). Default: all") String type
    ) {
        if ("alerting".equalsIgnoreCase(type)) {
            return prometheusService.getAlertingRules();
        } else if ("recording".equalsIgnoreCase(type)) {
            return prometheusService.getRecordingRules();
        }
        // Default: return both
        StringBuilder sb = new StringBuilder();
        sb.append("=== Alerting Rules ===\n");
        sb.append(prometheusService.getAlertingRules()).append("\n");
        sb.append("=== Recording Rules ===\n");
        sb.append(prometheusService.getRecordingRules());
        return sb.toString();
    }

    // =========================================================================
    // Status Tools
    // =========================================================================

    @Tool(description = "Get Prometheus server status and build information. "
            + "Shows version, build date, Go version, and runtime configuration. "
            + "Use this to: verify Prometheus is running, check version for compatibility, "
            + "get server info for troubleshooting, confirm configuration.")
    public String getPrometheusStatus() {
        return prometheusService.getServerStatus();
    }
}
