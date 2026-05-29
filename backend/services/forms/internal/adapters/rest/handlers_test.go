package rest

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"

	"sundance/backend/services/forms/internal/core"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

func newTestHandlers(services *ports.Services) *handlers {
	var buf bytes.Buffer

	app := &core.Application{
		Services: services,
		Logger:   slog.New(slog.NewTextHandler(&buf, nil)),
	}

	return newHandlers(app)
}

func Test_isBadRequest(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			"should yield true when err is ErrVersionLocked",
			domain.ErrVersionLocked,
			true,
		},
		{
			"should yield true when err is ErrInvalidVersion",
			domain.ErrInvalidVersion,
			true,
		},
		{
			"should yield true when err is ErrInvalidVersionStatus",
			domain.ErrInvalidVersionStatus,
			true,
		},
		{
			"should yield true when err is ErrDuplicateVersion",
			domain.ErrDuplicateVersion,
			true,
		},
		{
			"should yield true when err is ErrInvalidPosition",
			domain.ErrInvalidPosition,
			true,
		},
		{
			"should yield true when err is ErrDuplicatePosition",
			domain.ErrDuplicatePosition,
			true,
		},
		{
			"should yield true when err is ErrInvalidRuleType",
			domain.ErrInvalidRuleType,
			true,
		},
		{
			"should yield true when err is ErrDuplicateRuleType",
			domain.ErrDuplicateRuleType,
			true,
		},
		{
			"should yield true when err is ErrPublishedByRequired",
			domain.ErrPublishedByRequired,
			true,
		},
		{
			"should yield true when err is ErrRetiredByRequired",
			domain.ErrRetiredByRequired,
			true,
		},
		{
			"should yield true when err is ErrInvalidFieldType",
			domain.ErrInvalidFieldType,
			true,
		},
		{
			"should yield true when err is ErrInvalidFieldAttributes",
			domain.ErrInvalidFieldAttributes,
			true,
		},
		{
			"should yield true when err is ErrInvalidForm",
			domain.ErrInvalidForm,
			true,
		},
		{
			"should yield true when err is ErrFormHasActiveVersion",
			domain.ErrFormHasActiveVersion,
			true,
		},
		{
			"should yield true when err is ErrInvalidPage",
			domain.ErrInvalidPage,
			true,
		},
		{
			"should yield true when err is ErrInvalidSection",
			domain.ErrInvalidSection,
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
