package domain

type Lookup struct {
	Value any
	Label string
}

func NewLookup(value any, label string) *Lookup {
	return &Lookup{
		Value: value,
		Label: label,
	}
}
