package store

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/pablotrinidad/weatherreport/store/openweather"
)

// fixedWeatherResponse returned by mock API. Represents MEX airport (CDMX).
var fixedWeatherResponse openweather.WeatherItem = openweather.WeatherItem{
	Lat:             19.4360762,
	Lon:             -99.074097,
	Description:     []string{"cloudy", "foggy"},
	CityName:        "Mountain View",
	ObservationTime: 1601438975,
	Temp:            13,
	MaxTemp:         15,
	MinTemp:         10,
	FeelsLike:       14,
	Humidity:        60,
}

var fixedWeatherReport WeatherReport = WeatherReport{
	Lat:             19.4360762,
	Lon:             -99.074097,
	Description:     []string{"cloudy", "foggy"},
	CityName:        "Mountain View",
	Temp:            13,
	MaxTemp:         15,
	MinTemp:         10,
	FeelsLike:       14,
	Humidity:        60,
	ObservationTime: time.Unix(1601438975, 0),
}

var airports map[string]Airport = map[string]Airport{
	"TLC": {Code: "TLC", Latitude: 19.3371, Longitude: -99.566},
	"MTY": {Code: "MTY", Latitude: 25.7785, Longitude: -100.107},
	"MEX": {Code: "MEX", Latitude: 19.4363, Longitude: -99.0721},
	"TAM": {Code: "TAM", Latitude: 22.2964, Longitude: -97.8659},
}

func TestConcurrentStore_GetWeatherReport(t *testing.T) {
	tests := []struct {
		name         string
		queries      []Airport
		wantErrors   map[string]bool
		wantSuccess  map[string]bool
		apiMustFail  bool
		wantAPICalls int
	}{
		{
			name:         "empty airport list",
			queries:      []Airport{},
			wantAPICalls: 0,
			wantErrors:   map[string]bool{},
			wantSuccess:  map[string]bool{},
		},
		{
			name:         "single-element airport list",
			queries:      []Airport{airports["TLC"]},
			wantAPICalls: 1,
			wantErrors:   map[string]bool{},
			wantSuccess:  map[string]bool{"TLC": true},
		},
		{
			name:         "multiple unique airports",
			queries:      []Airport{airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"]},
			wantAPICalls: 4,
			wantErrors:   map[string]bool{},
			wantSuccess:  map[string]bool{"TLC": true, "MTY": true, "MEX": true, "TAM": true},
		},
		{
			name: "multiple repeated airports",
			queries: []Airport{
				airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
				airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
				airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
				airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
				airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
				airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"],
			},
			wantAPICalls: 4,
			wantErrors:   map[string]bool{},
			wantSuccess:  map[string]bool{"TLC": true, "MTY": true, "MEX": true, "TAM": true},
		},
		{
			name:         "failed API call",
			queries:      []Airport{airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"]},
			apiMustFail:  true,
			wantAPICalls: 4,
			wantErrors:   map[string]bool{"TLC": true, "MTY": true, "MEX": true, "TAM": true},
			wantSuccess:  map[string]bool{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := openweather.NewAPIMockClient(fixedWeatherResponse)
			api.FailNext = test.apiMustFail
			store := NewConcurrentStore(api)
			gotRes := store.GetWeatherByAirportCode(test.queries)
			if test.wantAPICalls != api.APICalls["coords"] {
				t.Errorf("GetWeatherByAirportCode(%v)\n called %d times Open Weather API, want %d calls", test.queries, api.APICalls["coords"], test.wantAPICalls)
			}
			if len(gotRes) != len(test.wantSuccess)+len(test.wantErrors) {
				t.Fatalf("GetWeatherByAirportCode(%v)\n returned %d results, want %d\nResponse:\n%v", test.queries, len(gotRes), len(test.wantSuccess)+len(test.wantErrors), gotRes)
			}
			for k, gotRes := range gotRes {
				wantRes := fixedWeatherReport
				switch gotRes.Failed {
				case true:
					if _, ok := test.wantErrors[k]; !ok {
						t.Errorf("want error for aiport code %q, got %v", k, gotRes)
						continue
					}
					wantRes.Failed = true
					wantRes.FailMessage = gotRes.FailMessage
					delete(test.wantErrors, k)
				case false:
					if _, ok := test.wantSuccess[k]; !ok {
						t.Errorf("want successful response for aiport code %q, got %v", k, gotRes)
						continue
					}
					delete(test.wantSuccess, k)
				}
				if diff := cmp.Diff(gotRes, wantRes); diff != "" {
					t.Errorf("got %v, want %v\ndiff: got->want %s", gotRes, wantRes, diff)
				}
			}
		})
	}
}
