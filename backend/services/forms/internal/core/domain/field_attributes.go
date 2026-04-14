package domain

type FieldAttributes interface{}

type BaseFieldAttributes struct {
	IsReadOnly bool
	IsRequired bool
}

type TextFieldAttributes struct {
	BaseFieldAttributes
	MinLength   *int
	MaxLength   *int
	Pattern     string
	Placeholder string
}

type NumberFieldAttributes struct {
	BaseFieldAttributes
	Min  *float64
	Max  *float64
	Step *float64
}

type SelectFieldAttributes struct {
	BaseFieldAttributes
	Multiple    bool
	MinSelected *int
	MaxSelected *int
}

type CheckboxFieldAttributes struct {
	BaseFieldAttributes
	IsCheckedByDefault bool
}

type DateFieldAttributes struct {
	BaseFieldAttributes
	MinDate *string
	MaxDate *string
}

