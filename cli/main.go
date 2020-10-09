package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pablotrinidad/weatherreport/store"
	"github.com/pablotrinidad/weatherreport/store/openweather"
)

type datasetFormat uint

const (
	unknownDatasetFormat datasetFormat = iota
	airportDatasetFormat
	citiesDatasetFormat
)

func init() {
	log.SetPrefix("")
	log.SetFlags(0)
}

func main() {
	deps, err := getApplicationDependencies()
	if err != nil {
		log.Fatalf("%v", err)
	}
	app := NewApp(deps)

	dataset, format, err := read()
	if err != nil {
		log.Fatalf("%v\nuse -h flag for usage instructions", err)
	}

	var report map[string]store.WeatherReport
	switch format {
	case airportDatasetFormat:
		airports, err := app.LoadAirportsDataset(dataset)
		if err != nil {
			log.Fatalf("Failed loading dataset:\n\t%v", err)
		}
		report, err = app.GetAirportsWeather(airports)
		if err != nil {
			log.Fatalf("Failed obtaining weather report:\n\t%v", err)
		}
	case citiesDatasetFormat:
		cities, err := app.LoadCitiesDataset(dataset)
		if err != nil {
			log.Fatalf("Failed loading dataset:\n\t%v\n", err)
		}
		report, err = app.GetCitiesWeather(cities)
		if err != nil {
			log.Fatalf("Failed obtaining weather report:\n\t%v", err)
		}
	}

	printResults(report)
}

// getApplicationDependencies returns newly initialized application dependencies.
func getApplicationDependencies() (*Deps, error) {
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("failed obtaining configuration: %v", err)
	}
	ow, err := openweather.NewAPIClient(config.openweatherAPIKey, "metric")
	if err != nil {
		return nil, fmt.Errorf("failed initializing OpenWeather API Client: %v", err)
	}
	return &Deps{store: store.NewConcurrentStore(ow)}, nil
}

// read command line flags.
func read() (string, datasetFormat, error) {
	var dataset string
	var format uint
	flag.StringVar(&dataset, "d", "", "path to dataset location")
	flag.UintVar(&format, "f", 0, "dataset format [1,2]:\n\t1: Airport codes dataset\n\t2: City names dataset")
	flag.Parse()
	if dataset == "" {
		return "", unknownDatasetFormat, fmt.Errorf("cannot use empty dataset location")
	}
	if format <= 0 || format > 2 {
		return "", unknownDatasetFormat, fmt.Errorf("got invalid dataset format %d, use 1 for airport codes dataset and 2 for city names dataset", format)
	}

	return dataset, datasetFormat(format), nil
}

// printResults upon confirmation.
func printResults(results map[string]store.WeatherReport) {
	fmt.Printf("\nDo you want to print %d results? [y/N]: ", len(results))
	if !confirmation() {
		fmt.Println("\nBYE 👋!")
		os.Exit(0)
	}
	for k, r := range results {
		fmt.Println("==========================================")
		fmt.Printf("q: %s\n", k)
		if r.Failed {
			fmt.Printf("\tcouldn't get weather information for %q\n", k)
			fmt.Printf("\treason: %s\n", r.FailMessage)
			continue
		}
		fmt.Printf("\tcity name: %s\n", r.CityName)
		fmt.Printf("\tlat:%0.2f lon: %0.2f\n", r.Lat, r.Lon)
		fmt.Printf("\tdescription: %v\n", r.Description)
		fmt.Printf("\ttemp: %0.2f°C\n", r.Temp)
		fmt.Printf("\t\tmax: %0.2f°C\n", r.MaxTemp)
		fmt.Printf("\t\tmin: %0.2f°C\n", r.MinTemp)
		fmt.Printf("\t\tfeels like: %0.2f°C\n", r.FeelsLike)
		fmt.Printf("\thumidity: %d%%\n", r.Humidity)
		fmt.Printf("\tobservation time (UTC): %v\n", r.ObservationTime)
	}
}

func confirmation() bool {
	for true {
		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			log.Fatal("failed reading choice from STDIN")
		}
		switch strings.ToLower(response) {
		case "y", "yes", "si", "sip":
			return true
		case "n", "no", "nel":
			return false
		default:
			fmt.Print("Please type (y)es or (n)o and then press enter: ")
		}
	}
	return false
}
