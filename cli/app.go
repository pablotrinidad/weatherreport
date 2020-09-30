package main

import (
	"fmt"

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

func (a *app) loadDataset() error {
	fmt.Println("LOADING DATASET...")
	return nil
}

func (a *app) MAGIC() error {
	fmt.Println("MAGIC STARTS")
	airports := []store.Airport{
		airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
		airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
		airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
		airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
		airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
		airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
		airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
		airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
	}
	data, err := a.deps.store.GetWeatherReport(airports)
	if err != nil {
		return err
	}
	for k, v := range data {
		fmt.Printf("\t%s\n", k)
		fmt.Printf("\t\tTemp: %.2f°C\n", v.Temp)
		fmt.Printf("\t\tFeels like: %.2f°C\n", v.FeelsLike)
		fmt.Printf("\t\tHumidity: %d\n", v.Humidity)
		fmt.Printf("\t\tObservation time: %s\n", v.ObservationTime)
	}
	fmt.Println("MAGIC ENDS")
	return nil
}
