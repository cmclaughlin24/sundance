package services

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
)

func Bootstrap(logger *slog.Logger, repository *ports.Repository) *ports.Services {
	return &ports.Services{
		Submissions: NewSubmissionsService(logger, repository),
	}
}
