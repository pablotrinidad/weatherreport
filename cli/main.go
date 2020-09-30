package main

import (
	"fmt"
	"os"

	"github.com/pablotrinidad/weatherreport/store"
	"github.com/pablotrinidad/weatherreport/store/openweather"
)

func main() {
	deps, err := getApplicationDependencies()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot start app:\n\t%v\n", err)
		os.Exit(1)
	}
	app := newApp(deps)

	if err := app.loadDataset(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed loading tickets dataset:\n\t%v\n", err)
		os.Exit(1)
	}

	if err := app.MAGIC(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot execute main operation:\n\t%v\n", err)
		os.Exit(1)
	}
}

func getApplicationDependencies() (*appDependencies, error) {
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("failed obtaining configuration: %v", err)
	}
	ow, err := openweather.NewAPIClient(config.openweatherAPIKey, "metric")
	if err != nil {
		return nil, fmt.Errorf("failed initializing OpenWeather API Client: %v", err)
	}
	return &appDependencies{store: store.NewConcurrentStore(ow)}, nil
}
