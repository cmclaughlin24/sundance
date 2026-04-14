package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/cmclaughlin24/sundance/common"
	"github.com/cmclaughlin24/sundance/tenants/internal/adapters/rest"
	"github.com/cmclaughlin24/sundance/tenants/internal/core"
)

const port = 80

func main() {
	settingsPath := flag.String("settings", "settings.json", "Path to settings JSON file")
	flag.Parse()

	var settings core.ApplicationSettings

	if err := common.ReadJsonFile(*settingsPath, &settings); err != nil {
		panic(err)
	}

	app, err := core.NewApplication(settings)

	if err != nil {
		panic(err)
	}

	defer app.Close()
	mux := rest.NewRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("application listening on :%d", port)

	if err := server.ListenAndServe(); err != nil {
		app.Logger.Fatal(err)
	}
}
