package tools

import (
	"strings"
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

func compileOperation(target spec.TargetConfig, path, method string, pathItem *openapi3.PathItem, op *openapi3.Operation) CompiledTool {
	var builder inputSchemaBuilder

	paramRefs := append([]*openapi3.ParameterRef{}, pathItem.Parameters...)
	paramRefs = append(paramRefs, op.Parameters...)

	for _, paramRef := range paramRefs {
		if paramRef == nil || paramRef.Value == nil {
			continue
		}

		param := paramRef.Value

		if strings.EqualFold("header", param.In) || strings.EqualFold("cookie", param.In) {
			continue
		}
	}

	return CompiledTool{}
}
