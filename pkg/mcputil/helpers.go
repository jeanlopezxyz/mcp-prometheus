package mcputil

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GetArguments extracts the tool call arguments from an MCP CallToolRequest.
func GetArguments(request *mcp.CallToolRequest) (map[string]any, error) {
	params, ok := request.GetParams().(*mcp.CallToolParamsRaw)
	if !ok {
		return nil, errors.New("invalid tool call parameters")
	}
	var args map[string]any
	if err := json.Unmarshal(params.Arguments, &args); err != nil {
		return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
	}
	return args, nil
}

// NewTextResult creates a successful text result.
func NewTextResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

// NewErrorResult creates an error text result.
func NewErrorResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}
