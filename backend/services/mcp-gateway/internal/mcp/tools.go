package mcp

import (
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type compiledTool struct {
	TargetName string
	BaseURL    string
	Method     string
	Path       string
	Headers    map[string]string
	Timeout    time.Duration
	Tool       mcp.Tool
}

type toolRegistry struct {
	tools map[string]compiledTool
}
