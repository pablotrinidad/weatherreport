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

func TestConcurrentStore_GetWeatherByAirportCode(t *testing.T) {
	tests := []struct {
		name        string
		queries     []Airport
		wantErrors  map[string]bool
		wantSuccess map[string]bool
		apiMustFail bool
		wantUsage   APIUsage
	}{
		{
			name:        "empty airport list",
			queries:     []Airport{},
			wantErrors:  map[string]bool{},
			wantSuccess: map[string]bool{},
			wantUsage:   APIUsage{SuccessfulCalls: 0, FailedCalls: 0},
		},
		{
			name:        "single-element airport list",
			queries:     []Airport{airports["TLC"]},
			wantErrors:  map[string]bool{},
			wantSuccess: map[string]bool{"TLC": true},
			wantUsage:   APIUsage{SuccessfulCalls: 1, FailedCalls: 0},
		},
		{
			name:        "multiple unique airports",
			queries:     []Airport{airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"]},
			wantErrors:  map[string]bool{},
			wantSuccess: map[string]bool{"TLC": true, "MTY": true, "MEX": true, "TAM": true},
			wantUsage:   APIUsage{SuccessfulCalls: 4, FailedCalls: 0},
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
			wantErrors:  map[string]bool{},
			wantSuccess: map[string]bool{"TLC": true, "MTY": true, "MEX": true, "TAM": true},
			wantUsage:   APIUsage{SuccessfulCalls: 4, FailedCalls: 0},
		},
		{
			name:        "failed API call",
			queries:     []Airport{airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"]},
			apiMustFail: true,
			wantUsage:   APIUsage{SuccessfulCalls: 0, FailedCalls: 4},
			wantErrors:  map[string]bool{"TLC": true, "MTY": true, "MEX": true, "TAM": true},
			wantSuccess: map[string]bool{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := openweather.NewAPIMockClient(fixedWeatherResponse)
			api.FailNext = test.apiMustFail
			store := NewConcurrentStore(api)
			gotRes := store.GetWeatherByAirportCode(test.queries)
			gotUsage := store.GetAPIUsage()
			if diff := cmp.Diff(gotUsage, test.wantUsage); diff != "" {
				t.Errorf("got usage %v, want %v\ndiff: got->want %s", gotUsage, test.wantUsage, diff)
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
				if diff := cmp.Diff(gotRes, wantRes); diff != "" && !test.apiMustFail {
					t.Errorf("got %v, want %v\ndiff: got->want %s", gotRes, wantRes, diff)
				}
			}
		})
	}
}

func TestConcurrentStore_GetWeatherByCityName(t *testing.T) {
	tests := []struct {
		name        string
		queries     []string
		wantErrors  map[string]bool
		wantSuccess map[string]bool
		apiMustFail bool
		wantUsage   APIUsage
	}{
		{
			name:        "empty cities list",
			queries:     []string{},
			wantErrors:  map[string]bool{},
			wantSuccess: map[string]bool{},
			wantUsage:   APIUsage{SuccessfulCalls: 0, FailedCalls: 0},
		},
		{
			name:        "single-element cities list",
			queries:     []string{"Mountain View"},
			wantErrors:  map[string]bool{},
			wantSuccess: map[string]bool{"Mountain View": true},
			wantUsage:   APIUsage{SuccessfulCalls: 1, FailedCalls: 0},
		},
		{
			name:        "multiple unique cities",
			queries:     []string{"New York City", "San Francisco", "Seattle", "Denver", "Houston"},
			wantErrors:  map[string]bool{},
			wantSuccess: map[string]bool{"New York City": true, "San Francisco": true, "Seattle": true, "Denver": true, "Houston": true},
			wantUsage:   APIUsage{SuccessfulCalls: 5, FailedCalls: 0},
		},
		{
			name: "multiple repeated cities",
			queries: []string{
				"New York City", "San Francisco", "Seattle", "Denver", "Houston",
				"New York City", "San Francisco", "Seattle", "Denver", "Houston",
				"New York City", "San Francisco", "Seattle", "Denver", "Houston",
				"New York City", "San Francisco", "Seattle", "Denver", "Houston",
			},
			wantErrors:  map[string]bool{},
			wantSuccess: map[string]bool{"New York City": true, "San Francisco": true, "Seattle": true, "Denver": true, "Houston": true},
			wantUsage:   APIUsage{SuccessfulCalls: 5, FailedCalls: 0},
		},
		{
			name:        "failed API call",
			queries:     []string{"New York City", "San Francisco", "Seattle", "Denver", "Houston"},
			apiMustFail: true,
			wantUsage:   APIUsage{SuccessfulCalls: 0, FailedCalls: 5},
			wantErrors:  map[string]bool{"New York City": true, "San Francisco": true, "Seattle": true, "Denver": true, "Houston": true},
			wantSuccess: map[string]bool{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := openweather.NewAPIMockClient(fixedWeatherResponse)
			api.FailNext = test.apiMustFail
			store := NewConcurrentStore(api)
			gotRes := store.GetWeatherByCityName(test.queries)
			gotUsage := store.GetAPIUsage()
			if diff := cmp.Diff(gotUsage, test.wantUsage); diff != "" {
				t.Errorf("got usage %v, want %v\ndiff: got->want %s", gotUsage, test.wantUsage, diff)
			}
			if len(gotRes) != len(test.wantSuccess)+len(test.wantErrors) {
				t.Fatalf("GetWeatherByCityName(%v)\n returned %d results, want %d\nResponse:\n%v", test.queries, len(gotRes), len(test.wantSuccess)+len(test.wantErrors), gotRes)
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
				if diff := cmp.Diff(gotRes, wantRes); diff != "" && !test.apiMustFail {
					t.Errorf("got %v, want %v\ndiff: got->want %s", gotRes, wantRes, diff)
				}
			}
		})
	}
}
