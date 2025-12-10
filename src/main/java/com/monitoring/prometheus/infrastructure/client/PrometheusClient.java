package com.monitoring.prometheus.infrastructure.client;

import com.monitoring.prometheus.infrastructure.dto.*;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import org.eclipse.microprofile.rest.client.annotation.RegisterClientHeaders;
import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;

@Path("/api/v1")
@RegisterRestClient(configKey = "prometheus-api")
@RegisterClientHeaders(KubernetesBearerTokenHeaderFactory.class)
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public interface PrometheusClient {

    /**
     * Execute an instant PromQL query.
     */
    @GET
    @Path("/query")
    QueryResponseDto query(
        @QueryParam("query") String query,
        @QueryParam("time") String time
    );

    /**
     * Execute a range PromQL query.
     */
    @GET
    @Path("/query_range")
    QueryResponseDto queryRange(
        @QueryParam("query") String query,
        @QueryParam("start") String start,
        @QueryParam("end") String end,
        @QueryParam("step") String step
    );

    /**
     * Get the current state of scrape targets.
     */
    @GET
    @Path("/targets")
    TargetsResponseDto getTargets(@QueryParam("state") String state);

    /**
     * Get alerting and recording rules.
     */
    @GET
    @Path("/rules")
    RulesResponseDto getRules(@QueryParam("type") String type);

    /**
     * Get active alerts.
     */
    @GET
    @Path("/alerts")
    AlertsResponseDto getAlerts();

    /**
     * Get Prometheus server metadata.
     */
    @GET
    @Path("/status/buildinfo")
    Object getBuildInfo();

    /**
     * Get Prometheus server runtime info.
     */
    @GET
    @Path("/status/runtimeinfo")
    Object getRuntimeInfo();

    /**
     * Get Prometheus server configuration.
     */
    @GET
    @Path("/status/config")
    Object getConfig();

    /**
     * Get label values for a specific label name.
     */
    @GET
    @Path("/label/{labelName}/values")
    Object getLabelValues(@PathParam("labelName") String labelName);

    /**
     * Get all label names.
     */
    @GET
    @Path("/labels")
    Object getLabels();
}
