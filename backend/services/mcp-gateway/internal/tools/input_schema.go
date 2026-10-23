package tools

const (
	propertyKeyPath  = "path"
	propertyKeyQuery = "query"
	propertyKeyBody  = "body"
	schemaTypeObject = "object"
)

type inputSchema struct {
	Type        string                 `json:"type,omitempty"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]inputSchema `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
}

type inputSchemaBuilder struct {
	path          map[string]inputSchema
	pathRequired  []string
	query         map[string]inputSchema
	queryRequired []string
	body          inputSchema
	bodyRequired  bool
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
