package inmemory

import (
	"log"

	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
)

func Bootstrap(logger *log.Logger) *ports.Repository {
	return &ports.Repository{
		Submissions: NewInMemorySubmissionsRepository(logger),
	}
}
