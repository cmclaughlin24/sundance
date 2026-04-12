package domain

type SectionID string

type Section struct {
	ID         SectionID
	Key        string
	Name       string
	Position   int
	Fields     map[int]*Field
	Conditions []*ConditionalRule
}

func NewSection(id SectionID, key, name string, position int) (*Section, error) {
	s := &Section{
		ID:       id,
		Key:      key,
		Name:     name,
		Position: position,
		Fields:   make(map[int]*Field),
	}

	// TODO: Implement domain specific validation.

	return s, nil
}
