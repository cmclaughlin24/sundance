package tools

import (
	"errors"
	"sundance/backend/services/mcp-gateway/internal/spec"

	"github.com/getkin/kin-openapi/openapi3"
)

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
