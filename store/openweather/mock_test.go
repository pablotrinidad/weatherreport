package openweather

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var fixedWeatherItem WeatherItem = WeatherItem{
	Lat:             1,
	Lon:             2,
	Temp:            3,
	FeelsLike:       4,
	Humidity:        5,
	ObservationTime: 1601438069,
}

func TestAPIMockClient_GetWeatherByCoords(t *testing.T) {
	tests := []struct {
		name     string
		failNext bool
		wantErr  bool
	}{
		{
			name: "base test",
		},
		{
			name:     "throws error",
			failNext: true,
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewAPIMockClient(fixedWeatherItem)
			c.FailNext = test.failNext
			gotRes, gotErr := c.GetWeatherByCoords(1, 2)
			if gotErr == nil && test.wantErr {
				t.Fatalf("GetWeatherByCoords(1, 2) returned nill error, want error")
			} else if gotErr != nil && !test.wantErr {
				t.Fatalf("GetWeatherByCoords(1, 2) returned unxepected error: %v", gotErr)
			}
			if !test.wantErr {
				if diff := cmp.Diff(*gotRes, fixedWeatherItem); diff != "" {
					t.Errorf("GetWeatherByCoords(1, 2): %v, want %v\ngot -> want diff: %s", gotRes, fixedWeatherItem, diff)
				}
			}
		})
	}
}

func TestAPIMockClient_GetWeatherByCityName(t *testing.T) {
	tests := []struct {
		name     string
		failNext bool
		wantErr  bool
	}{
		{
			name: "base test",
		},
		{
			name:     "throws error",
			failNext: true,
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewAPIMockClient(fixedWeatherItem)
			c.FailNext = test.failNext
			gotRes, gotErr := c.GetWeatherByCityName("Mountain view")
			if gotErr == nil && test.wantErr {
				t.Fatalf("GetWeatherByCityName('Mountain View') returned nill error, want error")
			} else if gotErr != nil && !test.wantErr {
				t.Fatalf("GetWeatherByCityName('Mountain View') returned unxepected error: %v", gotErr)
			}
			if !test.wantErr {
				if diff := cmp.Diff(*gotRes, fixedWeatherItem); diff != "" {
					t.Errorf("GetWeatherByCityName('Mountain View'): %v, want %v\ngot -> want diff: %s", gotRes, fixedWeatherItem, diff)
				}
			}
		})
	}
}
