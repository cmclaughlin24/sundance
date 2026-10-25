package tools

import (
	"sundance/backend/services/mcp-gateway/internal/spec"

	"github.com/getkin/kin-openapi/openapi3"
)

func ExtractTools(target spec.TargetConfig, doc *openapi3.T) []CompiledTool {
	var compiled []CompiledTool

	if doc == nil || doc.Paths == nil {
		return compiled
	}

	for path, pathItem := range doc.Paths.Map() {
		if pathItem == nil {
			continue
		}

		for method, op := range pathItem.Operations() {
			if op == nil {
				continue
			}

			if !target.Policy.IsAllowed(method, path) {
				continue
			}

			compiled = append(compiled, compileOperation(target, path, method, pathItem, op))
		}
	}

	return compiled
}
