FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o mcp-prometheus ./cmd/mcp-prometheus

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

LABEL io.modelcontextprotocol.server.name="io.github.jeanlopezxyz/mcp-prometheus"
LABEL io.k8s.display-name="MCP Prometheus Server"
LABEL io.openshift.tags="mcp,prometheus,monitoring,metrics"
LABEL maintainer="Jean Lopez"

WORKDIR /app
COPY --from=builder /app/mcp-prometheus /app/mcp-prometheus

USER 65532:65532
ENTRYPOINT ["/app/mcp-prometheus"]
EXPOSE 8080
