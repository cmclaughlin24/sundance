package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/adapters/persistence"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/adapters/rest"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/services"
)

type settings struct {
	Port        int                             `json:"port"`
	Persistence persistence.PersistenceSettings `json:"persistence"`
}

func main() {
	settingsPath := flag.String("settings", "settings.json", "Path to settings JSON file")
	flag.Parse()

	var settings settings

	if err := common.ReadJSONFile(*settingsPath, &settings); err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	r, err := persistence.Bootstrap(settings.Persistence, logger)

	if err != nil {
		logger.Fatal(err)
	}

	s := services.Bootstrap(logger, r)
	app := core.NewApplication(logger, r, s)

	defer app.Close()
	mux := rest.NewRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", settings.Port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("application listening on :%d", settings.Port)

	if err := server.ListenAndServe(); err != nil {
		app.Logger.Fatal(err)
	}
}
