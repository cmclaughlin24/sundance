package domain

type PageID string

type Page struct {
	ID         PageID
	Key        string
	Name       string
	Position   int
	Sections   map[int]*Section
	Conditions []*ConditionalRule
}

func NewPage(id PageID, key, name string, position int) (*Page, error) {
	p := &Page{
		ID:       id,
		Key:      key,
		Name:     name,
		Position: position,
		Sections: make(map[int]*Section),
	}

	// TODO: Implement domain specific validation.

	return p, nil
}
