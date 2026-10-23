package tools

import (
	"errors"
	"sundance/backend/services/mcp-gateway/internal/spec"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CompiledTool struct {
	TargetName string
	BaseURL    string
	Method     string
	Path       string
	Headers    map[string]string
	Timeout    time.Duration
	Tool       mcp.Tool
}

type Registry map[string]CompiledTool

func (r Registry) Register(tools ...CompiledTool) {
	for _, t := range tools {
		r[t.Tool.Name] = t
	}
}

func (r Registry) Get(name string) (CompiledTool, bool) {
	t, ok := r[name]
	return t, ok
}

func (r Registry) All() []CompiledTool {
	tools := make([]CompiledTool, 0, len(r))
	for _, t := range r {
		tools = append(tools, t)
	}

	return tools
}

func ExtractTools(target spec.TargetConfig, doc *openapi3.T) ([]CompiledTool, error) {
	var compiled []CompiledTool

	if doc == nil || doc.Paths == nil {
		return compiled, nil
	}

	for path, pathItem := range doc.Paths.Map() {
		if pathItem == nil {
			continue
		}

		for method, op := range pathItem.Operations() {
			if op == nil {
				continue
			}

			compiled = append(compiled, compileOperation(target, path, method, pathItem, op))
		}
	}

	return nil, errors.New("method not implemented")
}

func compileOperation(target spec.TargetConfig, path, method string, pathItem *openapi3.PathItem, op *openapi3.Operation) CompiledTool {
	return CompiledTool{}
}
