package domain

type FieldID string

type Field struct {
	ID         FieldID
	Key        string
	Name       string
	FieldType  string // TODO: Implement a field type
	Attributes FieldAttributes
	Position   int
	Conditions []*ConditionalRule
}

func NewField(id FieldID, key, name string, fieldType string, attributes FieldAttributes, position int) (*Field, error) {
	f := &Field{
		ID:         id,
		Key:        key,
		Name:       name,
		FieldType:  fieldType,
		Attributes: attributes,
		Position:   position,
	}

	// TODO: Implement domain specific validation.

	return f, nil
}
