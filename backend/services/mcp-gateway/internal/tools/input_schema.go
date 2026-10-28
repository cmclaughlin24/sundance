package tools

const (
	propertyKeyHeader = "header"
	propertyKeyPath   = "path"
	propertyKeyQuery  = "query"
	propertyKeyBody   = "body"
	schemaTypeObject  = "object"
	schemaTypeString  = "string"
)

type inputSchema struct {
	Type        string                 `json:"type,omitempty"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]inputSchema `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
	Items       *inputSchema           `json:"items,omitempty"`
	Format      string                 `json:"format,omitempty"`
	Pattern     string                 `json:"pattern,omitempty"`
	Enum        []any                  `json:"enum,omitempty"`
	Default     any                    `json:"default,omitempty"`
	Minimum     *float64               `json:"minimum,omitempty"`
	Maximum     *float64               `json:"maximum,omitempty"`
	MinLength   *uint64                `json:"minLength,omitempty"`
	MaxLength   *uint64                `json:"maxLength,omitempty"`
}

type inputSchemaBuilder struct {
	headers         map[string]inputSchema
	headersRequired []string
	path            map[string]inputSchema
	pathRequired    []string
	query           map[string]inputSchema
	queryRequired   []string
	body            inputSchema
	bodyRequired    bool
}

func (b *inputSchemaBuilder) addHeader(name string, s inputSchema, required bool) {
	if b.headers == nil {
		b.headers = make(map[string]inputSchema)
	}

	b.headers[name] = s

	if required {
		b.headersRequired = append(b.headersRequired, name)
	}
}

func (b *inputSchemaBuilder) addPath(name string, s inputSchema, required bool) {
	if b.path == nil {
		b.path = make(map[string]inputSchema)
	}

	b.path[name] = s

	if required {
		b.pathRequired = append(b.pathRequired, name)
	}
}

func (b *inputSchemaBuilder) addQuery(name string, s inputSchema, required bool) {
	if b.query == nil {
		b.query = make(map[string]inputSchema)
	}

	b.query[name] = s

	if required {
		b.queryRequired = append(b.queryRequired, name)
	}
}

func (b *inputSchemaBuilder) setBody(body inputSchema, required bool) {
	b.body = body
	b.bodyRequired = required
}

func (b *inputSchemaBuilder) build() inputSchema {
	properties := make(map[string]inputSchema)
	var required []string

	if len(b.headers) > 0 {
		properties[propertyKeyHeader] = inputSchema{
			Type:       schemaTypeObject,
			Properties: b.headers,
			Required:   b.headersRequired,
		}

		if len(b.headersRequired) > 0 {
			required = append(required, propertyKeyPath)
		}
	}

	if len(b.path) > 0 {
		properties[propertyKeyPath] = inputSchema{
			Type:       schemaTypeObject,
			Properties: b.path,
			Required:   b.pathRequired,
		}

		if len(b.pathRequired) > 0 {
			required = append(required, propertyKeyPath)
		}
	}

	if len(b.query) > 0 {
		properties[propertyKeyQuery] = inputSchema{
			Type:       schemaTypeObject,
			Properties: b.query,
			Required:   b.queryRequired,
		}

		if len(b.queryRequired) > 0 {
			required = append(required, propertyKeyQuery)
		}
	}

	if b.body.Type != "" || len(b.body.Properties) > 0 {
		properties[propertyKeyBody] = b.body

		if b.bodyRequired {
			required = append(required, propertyKeyBody)
		}
	}

	return inputSchema{
		Type:       schemaTypeObject,
		Properties: properties,
		Required:   required,
	}
}
