package inmemory

import (
	"log"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/ports"
)

func Bootstrap(logger *log.Logger) *ports.Repository {
	return &ports.Repository{}
}
