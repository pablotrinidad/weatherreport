package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pablotrinidad/weatherreport/store"
)

// Deps are an application dependencies.
type Deps struct {
	store store.Store
}

// App provides methods for reading datasets and performing weather queries.
type App struct {
	deps *Deps
}

// NewApp using dependencies/
func NewApp(deps *Deps) *App {
	return &App{deps: deps}
}

// LoadAirportsDataset from source file and returns full list of found airports (including duplicates).
func (a *App) LoadAirportsDataset(src string) ([]store.Airport, error) {
	log.Printf("loading airports from %s", src)
	rows, err := loadCSV(src)
	if err != nil {
		return nil, err
	}
	airports := make([]store.Airport, (len(rows))*2)
	unique := make(map[string]bool)
	for i, row := range rows {
		if len(row) < 6 {
			return nil, fmt.Errorf("missing columns at row %d, got %d", i+2, len(row))
		}
		coords := map[int]float64{2: 0.0, 3: 0.0, 4: 0.0, 5: 0.0}
		for k := range coords {
			c, err := strconv.ParseFloat(row[k], 64)
			if err != nil {
				var coordType string
				if k%2 == 0 {
					coordType = "lat"
				} else {
					coordType = "lon"
				}
				return nil, fmt.Errorf("got invalid coordinates (%s) value %q for airport code %q at row %d", coordType, row[k], row[(k/2)-1], i+2)
			}
			coords[k] = c
		}
		airports[i*2] = store.Airport{Code: strings.Trim(row[0], ""), Latitude: coords[2], Longitude: coords[3]}
		airports[i*2+1] = store.Airport{Code: strings.Trim(row[1], ""), Latitude: coords[4], Longitude: coords[5]}
		if airports[i*2].Code == "" {
			return nil, fmt.Errorf("got empty airport code (origin) at row %d", i+2)
		}
		if airports[i*2+1].Code == "" {
			return nil, fmt.Errorf("got empty airport code (destination) at row %d", i+2)
		}
		unique[airports[i*2].Code] = true
		unique[airports[i*2+1].Code] = true
	}
	log.Printf("\t✅  loaded %d airports (%d unique)", len(airports), len(unique))
	return airports, nil
}

// LoadCitiesDataset from source file and returns full list of found city names (including duplicates).
func (a *App) LoadCitiesDataset(src string) ([]string, error) {
	log.Printf("loading cities from %s", src)
	rows, err := loadCSV(src)
	if err != nil {
		return nil, err
	}
	cities := make([]string, len(rows))
	unique := make(map[string]bool)
	for i, row := range rows {
		if len(row) < 1 {
			return nil, fmt.Errorf("missing columns at row %d", i+2)
		}
		cities[i] = strings.Trim(row[0], " \n")
		if cities[i] == "" {
			return nil, fmt.Errorf("got empty city name at row %d", i+2)
		}
		unique[cities[i]] = true
	}
	log.Printf("\t✅  loaded %d cities (%d unique)", len(cities), len(unique))
	return cities, nil
}

func (a *App) GetAirportsWeather(airports []store.Airport) (map[string]store.WeatherReport, error) {
	log.Print("\nfetching weather information...")
	start := time.Now()
	results := a.deps.store.GetWeatherByAirportCode(airports)
	elapsed := time.Since(start)
	printReport(results, elapsed)
	return results, nil
}

func (a *App) GetCitiesWeather(cities []string) (map[string]store.WeatherReport, error) {
	log.Print("\nfetching weather information...")
	start := time.Now()
	results := a.deps.store.GetWeatherByCityName(cities)
	elapsed := time.Since(start)
	printReport(results, elapsed)
	return results, nil
}

func printReport(results map[string]store.WeatherReport, elapsed time.Duration) {
	log.Printf("\t✅  DONE")
	log.Printf("\tresults: %d", len(results))
	log.Printf("\telapsed time: %s", elapsed)

	var success, failed uint
	for _, r := range results {
		if r.Failed {
			failed++
		} else {
			success++
		}
	}
	log.Printf("\tsucessful: %d", success)
	log.Printf("\tfailed: %d", failed)
}

func loadCSV(src string) ([][]string, error) {
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return data[1:], nil
}
