package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/pablotrinidad/weatherreport/store"
)

var airports map[string]store.Airport = map[string]store.Airport{
	"TLC": {Code: "TLC", Latitude: 19.3371, Longitude: -99.566},
	"MTY": {Code: "MTY", Latitude: 25.7785, Longitude: -100.107},
	"MEX": {Code: "MEX", Latitude: 19.4363, Longitude: -99.0721},
	"TAM": {Code: "TAM", Latitude: 22.2964, Longitude: -97.8659},
}

type appDependencies struct {
	store store.Store
}

type app struct {
	deps *appDependencies
}

func newApp(deps *appDependencies) *app {
	return &app{deps: deps}
}

func (a *app) loadDataset(src string) ([]store.Airport, error) {
	fmt.Println("LOADING DATASET...")
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)

	vals, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	airports := make([]store.Airport, (len(vals)-1)*2)
	for i, row := range vals[1:] {
		coords := map[int]float64{2: 0.0, 3: 0.0, 4: 0.0, 5: 0.0}
		for k := range coords {
			c, err := strconv.ParseFloat(row[k], 64)
			if err != nil {
				return nil, err
			}
			coords[k] = c
		}
		airports[i*2] = store.Airport{Code: row[0], Latitude: coords[2], Longitude: coords[3]}
		airports[i*2+1] = store.Airport{Code: row[1], Latitude: coords[4], Longitude: coords[5]}
	}

	return airports, nil
}

func (a *app) MAGIC(airports []store.Airport) error {
	fmt.Println("MAGIC STARTS")
	data, err := a.deps.store.GetWeatherReport(airports)
	if err != nil {
		return err
	}
	for k, v := range data {
		fmt.Printf("\t%s\n", k)
		fmt.Printf("\t\tTemp: %.2f°C\n", v.Temp)
		fmt.Printf("\t\tFeels like: %.2f°C\n", v.FeelsLike)
		fmt.Printf("\t\tHumidity: %d%%\n", v.Humidity)
		fmt.Printf("\t\tObservation time: %s\n", v.ObservationTime)
	}
	fmt.Println(len(data))
	fmt.Println("MAGIC ENDS")
	return nil
}
