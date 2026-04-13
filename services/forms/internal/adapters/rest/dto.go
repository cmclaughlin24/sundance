package rest

type upsertFormDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenantId"`
}

