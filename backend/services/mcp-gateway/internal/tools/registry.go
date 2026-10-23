package tools

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
