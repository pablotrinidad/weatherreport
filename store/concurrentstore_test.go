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
	Temp:            13,
	FeelsLike:       12,
	Humidity:        69,
	ObservationTime: 1601438975,
}

var fixedWeatherReport WeatherReport = WeatherReport{
	Lat:             19.4360762,
	Lon:             -99.074097,
	Temp:            13,
	FeelsLike:       12,
	Humidity:        69,
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
		name            string
		queries         []Airport
		apiMustFail     bool
		wantMinAPICalls uint
		wantErr         bool
		wantRes         map[string]WeatherReport
	}{
		{
			name:            "empty airport list",
			queries:         []Airport{},
			wantMinAPICalls: 0,
			wantErr:         false,
			wantRes:         map[string]WeatherReport{},
		},
		{
			name:            "single-element airport list",
			queries:         []Airport{airports["TLC"]},
			wantMinAPICalls: 1,
			wantErr:         false,
			wantRes:         map[string]WeatherReport{"TLC": fixedWeatherReport},
		},
		{
			name:            "multiple unique airports",
			queries:         []Airport{airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"]},
			wantMinAPICalls: 4,
			wantErr:         false,
			wantRes: map[string]WeatherReport{
				"TLC": fixedWeatherReport,
				"MTY": fixedWeatherReport,
				"MEX": fixedWeatherReport,
				"TAM": fixedWeatherReport,
			},
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
			wantMinAPICalls: 4,
			wantErr:         false,
			wantRes: map[string]WeatherReport{
				"TLC": fixedWeatherReport,
				"MTY": fixedWeatherReport,
				"MEX": fixedWeatherReport,
				"TAM": fixedWeatherReport,
			},
		},
		{
			name:            "failed API call",
			queries:         []Airport{airports["TLC"], airports["MTY"], airports["MEX"], airports["TAM"]},
			apiMustFail:     true,
			wantMinAPICalls: 1,
			wantErr:         true,
			wantRes:         nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := openweather.NewAPIMockClient(fixedWeatherResponse)
			api.FailNext = test.apiMustFail
			store := NewConcurrentStore(api)
			gotRes, gotErr := store.GetWeatherReport(test.queries)
			if gotErr == nil && test.wantErr {
				t.Fatalf("GetWeatherReport(%v)\n returned nil error, want error", test.queries)
			}
			if gotErr != nil && !test.wantErr {
				t.Fatalf("GetWeatherReport(%v)\n returned unexpexted error: %v", test.queries, gotErr)
			}
			if test.wantMinAPICalls > api.GetCurrentWeatherCalls {
				t.Errorf("GetWeatherReport(%v)\n called %d times Open Weather API, want %d calls at most", test.queries, api.GetCurrentWeatherCalls, test.wantMinAPICalls)
			}
			if diff := cmp.Diff(gotRes, test.wantRes); !test.wantErr && diff != "" {
				t.Errorf("GetWeatherReport(%v)\n: %v, want %v;\ndiff got -> want:\n %s ", test.queries, gotRes, test.wantRes, diff)
			}
		})
	}
}
