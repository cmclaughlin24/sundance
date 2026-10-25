package mcp

import (
	"encoding/json"
	"fmt"
	"sundance/backend/services/mcp-gateway/internal/tools"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func readRequestArgs(ctr *mcp.CallToolRequest) (*tools.Args, error) {
	var args tools.Args

	if len(ctr.Params.Arguments) == 0 {
		return &args, nil
	}

	if err := json.Unmarshal(ctr.Params.Arguments, &args); err != nil {
		return nil, fmt.Errorf("invalid arguments structure: %w", err)
	}

	return &args, nil
}
