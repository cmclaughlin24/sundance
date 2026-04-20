package domain

type FieldID string

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeSelect   FieldType = "select"
	FieldTypeCheckbox FieldType = "checkbox"
	FieldTypeDate     FieldType = "date"
)

type Field struct {
	ID         FieldID
	Key        string
	Name       string
	Type       FieldType
	Attributes FieldAttributes
	Position   int
	baseWithRules
}

func NewField(id FieldID, key, name string, fieldType FieldType, attributes FieldAttributes, position int) (*Field, error) {
	f := &Field{
		ID:         id,
		Key:        key,
		Name:       name,
		Type:       fieldType,
		Attributes: attributes,
		Position:   position,
	}

	// TODO: Implement domain specific validation.

	return f, nil
}
