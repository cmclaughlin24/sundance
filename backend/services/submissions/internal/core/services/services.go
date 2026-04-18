package services

import (
	"log"

	"github.com/cmclaughlin24/sundance/submissions/internal/core/ports"
)

func Bootstrap(logger *log.Logger, repository *ports.Repository) *ports.Services {
	return &ports.Services{
		Submissions: NewSubmissionsService(logger, repository),
	}
}
