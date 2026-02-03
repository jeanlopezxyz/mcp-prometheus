# Example Prompts for Prometheus MCP Server

This document provides example prompts that demonstrate how AI assistants can use the Prometheus MCP Server tools effectively.

## Cluster Health and Diagnostics

### Check cluster health metrics
> "Check the overall health of my Kubernetes cluster using Prometheus metrics. Look at node status, pod restarts, and resource utilization."

### Investigate high CPU usage
> "My cluster nodes are showing high CPU usage. Query Prometheus to find which pods and namespaces are consuming the most CPU resources."

### Check memory pressure
> "Are any nodes in my cluster experiencing memory pressure? Show me the current memory utilization across all nodes."

## PromQL Queries

### Run an instant query
> "Run a PromQL query to show the current CPU usage rate across all pods: `rate(container_cpu_usage_seconds_total[5m])`"

### Run a range query
> "Query Prometheus for the memory usage trend over the last 24 hours for pods in the `production` namespace."

### Compare metrics over time ranges
> "Compare the request latency for my API service between today and yesterday. Show the p95 latency values."

## Resource Analysis

### Top resource consumers
> "Which pods are consuming the most memory in the cluster right now? Show me the top 10."

### Namespace resource usage
> "Show me a breakdown of CPU and memory usage by namespace across the cluster."

### Disk usage analysis
> "Check the persistent volume usage across the cluster. Are any volumes close to running out of space?"

## Alerting Context

### Check alert-related metrics
> "What Prometheus recording rules are configured for alerting? Show me the current values of any critical threshold metrics."

### Service availability
> "Check the availability metrics for services in the `production` namespace. Are there any services with low uptime?"

## Tips for Effective Prompts

- **Be specific about time ranges** - Prometheus stores time-series data, so specifying time ranges helps get more useful results.
- **Reference namespaces** - When working with Kubernetes, specifying the namespace narrows down results.
- **Use PromQL when you know it** - If you know the exact PromQL query, include it in your prompt for direct execution.
- **Ask for comparisons** - The range query capabilities allow comparing metrics across different time periods.
