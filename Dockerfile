# =============================================================================
# MCP Prometheus Server - Quarkus Native Multi-stage Build
# =============================================================================
# Builds a native executable using Mandrel (Red Hat's GraalVM distribution)
# and deploys on the smallest possible image (ubi9-quarkus-micro-image)
#
# Build:
#   docker build -t ghcr.io/jeanlopezxyz/mcp-prometheus .
#
# Run:
#   docker run -i --rm -p 9081:9081 -e PROMETHEUS_URL=http://prometheus:9090 ghcr.io/jeanlopezxyz/mcp-prometheus
#
# Image size: ~50-100MB (vs ~400MB with JVM)
# Startup time: ~50ms (vs ~2s with JVM)
# =============================================================================

# Stage 1: Build Native Executable
FROM quay.io/quarkus/ubi-quarkus-mandrel-builder-image:jdk-21 AS build

WORKDIR /build

# Copy all source files
COPY --chown=quarkus:quarkus mvnw .
COPY --chown=quarkus:quarkus .mvn .mvn
COPY --chown=quarkus:quarkus pom.xml .
COPY --chown=quarkus:quarkus src src

# Build native executable (single step to avoid permission issues)
USER quarkus
RUN ./mvnw package -DskipTests -Dnative -B

# Stage 2: Runtime (Micro Image - smallest possible)
FROM quay.io/quarkus/ubi9-quarkus-micro-image:2.0

LABEL maintainer="Jean Lopez"
LABEL description="MCP Server for Prometheus (Native)"
LABEL io.k8s.display-name="MCP Prometheus Server"
LABEL io.openshift.tags="mcp,prometheus,promql,monitoring,metrics,quarkus,native"

WORKDIR /work/

# Setup permissions
RUN chown 1001 /work \
    && chmod "g+rwX" /work \
    && chown 1001:root /work

# Copy native executable from build stage
COPY --from=build --chown=1001:root --chmod=0755 /build/target/*-runner /work/application

EXPOSE 9081

USER 1001

# Environment variables
ENV QUARKUS_HTTP_HOST=0.0.0.0
ENV QUARKUS_HTTP_PORT=9081
ENV PROMETHEUS_URL=http://prometheus:9090

ENTRYPOINT ["./application"]
