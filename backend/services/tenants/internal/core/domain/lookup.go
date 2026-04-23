package domain

type Lookup struct {
	Value string
	Label string
}

func NewLookup(value, label string) *Lookup {
	return &Lookup{
		Value: value,
		Label: label,
	}
}
