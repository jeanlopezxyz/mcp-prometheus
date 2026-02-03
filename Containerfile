FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o mcp-prometheus ./cmd/mcp-prometheus

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/mcp-prometheus /mcp-prometheus
USER 65532:65532
ENTRYPOINT ["/mcp-prometheus"]
