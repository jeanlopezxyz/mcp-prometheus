package com.monitoring.prometheus.application.service;

import com.monitoring.prometheus.infrastructure.client.PrometheusClient;
import com.monitoring.prometheus.infrastructure.dto.*;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.eclipse.microprofile.rest.client.inject.RestClient;
import org.jboss.logging.Logger;

import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@ApplicationScoped
public class PrometheusService {

    private static final Logger LOG = Logger.getLogger(PrometheusService.class);

    @Inject
    @RestClient
    PrometheusClient prometheusClient;

    /**
     * Execute an instant PromQL query.
     */
    public String query(String promql) {
        LOG.infof("Executing PromQL query: %s", promql);
        try {
            QueryResponseDto response = prometheusClient.query(promql, null);
            return formatQueryResponse(response);
        } catch (Exception e) {
            LOG.errorf("Error executing query: %s", e.getMessage());
            return "Error: " + e.getMessage();
        }
    }

    /**
     * Execute a range PromQL query.
     */
    public String queryRange(String promql, String duration, String step) {
        LOG.infof("Executing PromQL range query: %s, duration: %s, step: %s", promql, duration, step);
        try {
            Instant now = Instant.now();
            Instant start = parseDuration(now, duration);

            QueryResponseDto response = prometheusClient.queryRange(
                promql,
                start.toString(),
                now.toString(),
                step != null ? step : "1m"
            );
            return formatQueryResponse(response);
        } catch (Exception e) {
            LOG.errorf("Error executing range query: %s", e.getMessage());
            return "Error: " + e.getMessage();
        }
    }

    /**
     * Get scrape targets status.
     */
    public String getTargets(String state) {
        LOG.infof("Getting targets with state: %s", state);
        try {
            TargetsResponseDto response = prometheusClient.getTargets(state);
            return formatTargetsResponse(response);
        } catch (Exception e) {
            LOG.errorf("Error getting targets: %s", e.getMessage());
            return "Error: " + e.getMessage();
        }
    }

    /**
     * Get alerting rules.
     */
    public String getAlertingRules() {
        LOG.info("Getting alerting rules");
        try {
            RulesResponseDto response = prometheusClient.getRules("alert");
            return formatRulesResponse(response);
        } catch (Exception e) {
            LOG.errorf("Error getting alerting rules: %s", e.getMessage());
            return "Error: " + e.getMessage();
        }
    }

    /**
     * Get recording rules.
     */
    public String getRecordingRules() {
        LOG.info("Getting recording rules");
        try {
            RulesResponseDto response = prometheusClient.getRules("record");
            return formatRulesResponse(response);
        } catch (Exception e) {
            LOG.errorf("Error getting recording rules: %s", e.getMessage());
            return "Error: " + e.getMessage();
        }
    }

    /**
     * Get Prometheus server status.
     */
    public String getServerStatus() {
        LOG.info("Getting server status");
        try {
            Object buildInfo = prometheusClient.getBuildInfo();
            Object runtimeInfo = prometheusClient.getRuntimeInfo();
            return String.format("Build Info:\n%s\n\nRuntime Info:\n%s", buildInfo, runtimeInfo);
        } catch (Exception e) {
            LOG.errorf("Error getting server status: %s", e.getMessage());
            return "Error: " + e.getMessage();
        }
    }

    // =========================================================================
    // Private Helper Methods
    // =========================================================================

    private Instant parseDuration(Instant from, String duration) {
        if (duration == null || duration.isEmpty()) {
            return from.minus(1, ChronoUnit.HOURS);
        }

        String value = duration.replaceAll("[^0-9]", "");
        String unit = duration.replaceAll("[0-9]", "").toLowerCase();
        int amount = Integer.parseInt(value);

        return switch (unit) {
            case "m" -> from.minus(amount, ChronoUnit.MINUTES);
            case "h" -> from.minus(amount, ChronoUnit.HOURS);
            case "d" -> from.minus(amount, ChronoUnit.DAYS);
            default -> from.minus(1, ChronoUnit.HOURS);
        };
    }

    private String formatQueryResponse(QueryResponseDto response) {
        if (!"success".equals(response.status())) {
            return String.format("Query failed: %s - %s", response.errorType(), response.error());
        }

        if (response.data() == null || response.data().result() == null) {
            return "No data returned";
        }

        StringBuilder sb = new StringBuilder();
        sb.append("Result Type: ").append(response.data().resultType()).append("\n\n");

        for (QueryResponseDto.QueryResultDto result : response.data().result()) {
            sb.append("Metric: ").append(formatMetricLabels(result.metric())).append("\n");

            if (result.value() != null) {
                sb.append("Value: ").append(result.value().get(1)).append("\n");
            }

            if (result.values() != null && !result.values().isEmpty()) {
                sb.append("Values: ").append(result.values().size()).append(" samples\n");
                // Show last 5 values
                int start = Math.max(0, result.values().size() - 5);
                for (int i = start; i < result.values().size(); i++) {
                    List<Object> v = result.values().get(i);
                    sb.append("  ").append(v.get(0)).append(": ").append(v.get(1)).append("\n");
                }
            }
            sb.append("\n");
        }

        return sb.toString();
    }

    private String formatMetricLabels(Map<String, String> labels) {
        if (labels == null || labels.isEmpty()) {
            return "{}";
        }
        return labels.entrySet().stream()
            .map(e -> e.getKey() + "=\"" + e.getValue() + "\"")
            .collect(Collectors.joining(", ", "{", "}"));
    }

    private String formatTargetsResponse(TargetsResponseDto response) {
        if (!"success".equals(response.status())) {
            return "Failed to get targets";
        }

        StringBuilder sb = new StringBuilder();
        sb.append("=== Active Targets ===\n\n");

        if (response.data().activeTargets() != null) {
            for (var target : response.data().activeTargets()) {
                sb.append("Pool: ").append(target.scrapePool()).append("\n");
                sb.append("URL: ").append(target.scrapeUrl()).append("\n");
                sb.append("Health: ").append(target.health()).append("\n");
                if (target.lastError() != null && !target.lastError().isEmpty()) {
                    sb.append("Last Error: ").append(target.lastError()).append("\n");
                }
                sb.append("Labels: ").append(target.labels()).append("\n\n");
            }
        }

        if (response.data().droppedTargets() != null && !response.data().droppedTargets().isEmpty()) {
            sb.append("=== Dropped Targets: ").append(response.data().droppedTargets().size()).append(" ===\n");
        }

        return sb.toString();
    }

    private String formatRulesResponse(RulesResponseDto response) {
        if (!"success".equals(response.status())) {
            return "Failed to get rules";
        }

        StringBuilder sb = new StringBuilder();

        if (response.data().groups() != null) {
            for (var group : response.data().groups()) {
                sb.append("=== Group: ").append(group.name()).append(" ===\n");
                sb.append("File: ").append(group.file()).append("\n\n");

                if (group.rules() != null) {
                    for (var rule : group.rules()) {
                        sb.append("Rule: ").append(rule.name()).append("\n");
                        sb.append("Type: ").append(rule.type()).append("\n");
                        sb.append("Query: ").append(rule.query()).append("\n");
                        sb.append("State: ").append(rule.state()).append("\n");
                        sb.append("Health: ").append(rule.health()).append("\n");

                        if (rule.alerts() != null && !rule.alerts().isEmpty()) {
                            sb.append("Active Alerts: ").append(rule.alerts().size()).append("\n");
                        }
                        sb.append("\n");
                    }
                }
            }
        }

        return sb.toString();
    }
}
