# Project Agents.md for Prometheus MCP Server

This Agents.md file provides comprehensive guidance for AI assistants and coding agents (like Gemini, Cursor, and others) to work with this codebase.

This repository contains the mcp-prometheus project,
a Go-based Model Context Protocol (MCP) server that provides Prometheus monitoring integration.
This MCP server enables AI assistants to query Prometheus metrics, run PromQL queries, and perform cluster diagnostics using the Model Context Protocol (MCP).

## Project Structure and Repository Layout

- Go package layout follows standard Go conventions:
  - `cmd/mcp-prometheus/` - main application entry point.
  - `pkg/` - libraries grouped by domain.
    - `kubernetes/` - Kubernetes client management and auto-discovery.
    - `mcputil/` - MCP utility functions and helpers.
    - `prometheus/` - Prometheus client and query execution.
    - `toolsets/` - Toolset registration and management for MCP tools.
    - `version/` - Version information management.
- `.github/` - GitHub-related configuration (Actions workflows, Dependabot).
- `build/` - Build configuration files (included via Makefile).
- `npm/` - Node packages that wrap the compiled binaries for distribution through npmjs.com.
- `Containerfile` - Container image description file.
- `Makefile` - Tasks for building, formatting, and testing.

## Feature Development

Implement new functionality in the Go sources under `cmd/` and `pkg/`.
The JavaScript (`npm/`) directory only wraps the compiled binary for distribution (npm).
Most changes will not require touching it unless the version or packaging needs to be updated.

### Adding New MCP Tools

The project uses a toolset-based architecture for organizing MCP tools:

- **Tool definitions** and handlers are created in `pkg/toolsets/`.
- **Toolsets** group related tools together (e.g., query tools, diagnostic tools).
- Each tool has a handler function that implements the tool's logic.

When adding a new tool:
1. Define the tool handler function that implements the tool's logic.
2. Register the tool in the appropriate toolset in `pkg/toolsets/`.
3. If creating a new toolset, register it with the MCP server.

## Building

Use the provided Makefile targets:

```bash
# Build the binary
make build

# Build for all supported platforms
make build-all-platforms
```

The resulting executable is `mcp-prometheus`.

## Running

```bash
# Using npx (Node.js package runner)
npx -y mcp-prometheus@latest

# Using the MCP Inspector
make build
npx @modelcontextprotocol/inspector@latest $(pwd)/mcp-prometheus

# Binary execution
./mcp-prometheus
```

## Tests

Run all Go tests with:

```bash
make test
```

## Linting

```bash
make fmt   # Format code
make vet   # Vet code
```

## Additional Makefile Targets

```bash
make help               # Display all available targets
make tidy               # Tidy go modules
make clean              # Clean build artifacts
make build-all-platforms # Cross-compile for all platforms
```

## Dependencies

When introducing new modules run `make tidy` so that `go.mod` and `go.sum` remain tidy.

## Coding Style

- Go modules target Go **1.25** (see `go.mod`).
- Build and test steps are defined in the Makefile - keep them working.
- Follow standard Go conventions for naming, formatting, and error handling.

## Distribution Methods

- **Native binaries** for Linux, macOS, and Windows are available in the GitHub releases.
- An **npm** package is available at [npmjs.com](https://www.npmjs.com/package/mcp-prometheus).
  It wraps the platform-specific binary and provides a convenient way to run the server using `npx`.
- A **container image** can be built using the `Containerfile`.
