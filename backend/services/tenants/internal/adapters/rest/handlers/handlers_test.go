package handlers

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"sundance/backend/services/tenants/internal/adapters/rest/dto"
	"sundance/backend/services/tenants/internal/core"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

func newTestHandlers(services *ports.Services) *Handlers {
	var buf bytes.Buffer

	app := &core.Application{
		Services: services,
		Logger:   slog.New(slog.NewTextHandler(&buf, nil)),
	}

	return NewHandlers(app)
}

func Test_isBadRequest(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			"should yield true when err is ErrDataSourceAttrParse",
			dto.ErrDataSourceAttrParse,
			true,
		},
		{
			"should yield true when err is ErrInvalidSourceType",
			domain.ErrInvalidSourceType,
			true,
		},
		{
			"should yield true when err is ErrInvalidSourceTypeAttributes",
			domain.ErrInvalidSourceTypeAttributes,
			true,
		},
		{
			"should yield false otherwise",
			errors.New("unknown error"),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBadRequest(tt.err)

			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}
