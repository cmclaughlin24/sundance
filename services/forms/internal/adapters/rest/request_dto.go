package rest

type upsertFormDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type createVersionDto struct{}

type upsertPageDto struct{}

type upsertSectionDto struct{}

type upsertFieldDto struct{}
