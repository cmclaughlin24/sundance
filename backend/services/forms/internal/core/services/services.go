package services

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

func Bootstrap(logger *slog.Logger, repository *ports.Repository) *ports.Services {
	return &ports.Services{
		Forms: NewFormsService(logger, repository),
	}
}
