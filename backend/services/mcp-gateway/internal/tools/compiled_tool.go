package tools

import (
	"fmt"
	"net/http"
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
	Headers    spec.HeadersConfig
	Timeout    time.Duration
	Tool       mcp.Tool
}

func (ct *CompiledTool) setHeaders(req *http.Request) *http.Request {
	if len(ct.Headers) == 0 {
		return req
	}

	for k, v := range ct.Headers {
		req.Header.Set(k, fmt.Sprint(v.Value))
	}

	return req
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

		if strings.EqualFold("cookie", param.In) {
			continue
		}

		paramSchema := resolveSchema(param.Schema, param.Description)

		switch strings.ToLower(param.In) {
		case "header":
			if target.Headers.AllowHeader(param.Name) {
				builder.addHeader(param.Name, paramSchema, param.Required)
			}
		case "path":
			builder.addPath(param.Name, paramSchema, param.Required)
		case "query":
			builder.addQuery(param.Name, paramSchema, param.Required)
		}
	}

	if op.RequestBody != nil && op.RequestBody.Value != nil {
		if media := op.RequestBody.Value.Content.Get("application/json"); media != nil && media.Schema != nil {
			bodySchema := resolveSchema(media.Schema, "")
			builder.setBody(bodySchema, op.RequestBody.Value.Required)
		}
	}

	return CompiledTool{
		TargetName: target.Name,
		BaseURL:    target.BaseURL,
		Method:     strings.ToUpper(method),
		Path:       path,
		Headers:    target.Headers,
		Timeout:    target.Timeout,
		Tool: mcp.Tool{
			Name:        formatToolName(target.Name, method, path),
			Description: formToolDescription(op),
			InputSchema: builder.build(),
		},
	}
}

func resolveSchema(schemaRef *openapi3.SchemaRef, fallbackDescription string) inputSchema {
	if schemaRef == nil || schemaRef.Value == nil {
		return inputSchema{Type: schemaTypeString, Description: fallbackDescription}
	}

	schema := schemaRef.Value
	description := schema.Description
	var typeStr string

	if strings.Trim(description, " ") == "" {
		description = fallbackDescription
	}

	if schema.Type != nil && len(schema.Type.Slice()) > 0 {
		typeStr = schema.Type.Slice()[0]
	}

	return inputSchema{
		Type:        typeStr,
		Description: description,
		Properties:  resolveProperties(schema),
		Required:    schema.Required,
		Items:       resolveItems(schema),
		Format:      schema.Format,
		Pattern:     schema.Pattern,
		Enum:        schema.Enum,
		Default:     schema.Default,
		Minimum:     schema.Min,
		Maximum:     schema.Max,
		MinLength:   &schema.MinLength,
		MaxLength:   schema.MaxLength,
	}
}

func resolveProperties(schema *openapi3.Schema) map[string]inputSchema {
	var props map[string]inputSchema

	if len(schema.Properties) == 0 {
		return props
	}

	props = make(map[string]inputSchema, len(schema.Properties))
	for propName, propRef := range schema.Properties {
		props[propName] = resolveSchema(propRef, "")
	}

	return props
}

func resolveItems(schema *openapi3.Schema) *inputSchema {
	if schema.Items == nil {
		return nil
	}

	resolved := resolveSchema(schema.Items, "")
	return &resolved
}

func formatToolName(targetName, method, path string) string {
	cleanPath := strings.Trim(path, "/")
	cleanPath = strings.ReplaceAll(cleanPath, "/", "_")
	cleanPath = strings.ReplaceAll(cleanPath, "-", "_")
	cleanPath = strings.ReplaceAll(cleanPath, "{", "")
	cleanPath = strings.ReplaceAll(cleanPath, "}", "")
	snakeName := fmt.Sprintf("%s-%s_%s", targetName, strings.ToLower(method), cleanPath)

	return snakeToCamel(snakeName)
}

func formToolDescription(op *openapi3.Operation) string {
	var description []string

	if strings.Trim(op.Summary, " ") != "" {
		description = append(description, op.Summary)
	}

	if strings.Trim(op.Description, " ") != "" {
		description = append(description, op.Description)
	}

	return strings.Join(description, " - ")
}

func snakeToCamel(value string) string {
	parts := strings.Split(value, "_")
	for i := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(parts[i])
		}

		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
		}
	}

	return strings.Join(parts, "")
}
