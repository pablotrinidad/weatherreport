package openweather

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func newTestAPIClient(apiKey, units string, server *httptest.Server, malformedURL bool) API {
	c, _ := NewAPIClient(apiKey, units)
	c.client = server.Client()
	if malformedURL {
		c.apiURL = "i'm not a valid HTTP URL :D"
	} else {
		c.apiURL = server.URL
	}
	return c
}

func newTestServer(wantStatusCode int, wantRes []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(wantStatusCode)
		w.Write(wantRes)
	}))
}

func compareResults(t *testing.T, caller string, gotRes, wantRes *WeatherItem, gotErr error, wantErr bool) {
	t.Helper()
	if wantErr {
		if gotErr == nil {
			t.Fatalf("%s returned nil error, want error", caller)
		}
		return
	}
	if gotErr != nil {
		t.Fatalf("%s returned unexpected error: %v", caller, gotErr)
	}
	if diff := cmp.Diff(gotRes, wantRes); diff != "" {
		t.Errorf("%s: %v, want %v\ngot -> want diff: %s", caller, gotRes, wantRes, diff)
	}
}

func TestNewAPIClient(t *testing.T) {
	tests := []struct {
		name          string
		apiKey, units string
		wantError     bool
	}{
		{
			name:   "successful creation",
			apiKey: "a",
			units:  "metric",
		},
		{
			name:      "empty api key",
			apiKey:    "",
			units:     "metric",
			wantError: true,
		},
		{
			name:      "invalid units value",
			apiKey:    "a",
			units:     "american",
			wantError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewAPIClient(test.apiKey, test.units)
			if err != nil && !test.wantError {
				t.Fatalf("NewAPIClient(%s, %s) returned unexpected error: %v", test.apiKey, test.units, err)
			}
			if err == nil && test.wantError {
				t.Fatalf("NewAPIClient(%s, %s) returned nil error, want error", test.apiKey, test.units)
			}
			if !test.wantError && got.apiKey != test.apiKey {
				t.Errorf("NewAPIClient(%s, %s) returned client with API key %s, want %s", test.apiKey, test.units, got.apiKey, test.apiKey)
			}
			if !test.wantError && got.units != test.units {
				t.Errorf("NewAPIClient(%s, %s) returned client with units %s, want %s", test.apiKey, test.units, got.units, test.units)
			}
		})
	}
}

func TestAPIClient_OWCurrent(t *testing.T) {
	tests := []struct {
		name          string
		lat, lon      float64
		cityName      string
		apiRes        []byte
		apiStatusCode int
		malformedURL  bool
		wantRes       *WeatherItem
		wantErr       bool
		closeServer   bool
	}{
		{
			name:     "successful response",
			cityName: "Mountain View",
			lat:      37.39, lon: -122.08,
			apiStatusCode: http.StatusOK,
			apiRes: []byte(`{
				"base": "stations",
				"clouds": {
					"all": 90
				},
				"cod": 200,
				"coord": {
					"lat": 37.39,
					"lon": -122.08
				},
				"dt": 1601662295,
				"id": 5375480,
				"main": {
					"feels_like": 27.8,
					"humidity": 30,
					"pressure": 1016,
					"temp": 28.87,
					"temp_max": 31.67,
					"temp_min": 27
				},
				"name": "Mountain View",
				"sys": {
					"country": "US",
					"id": 5845,
					"sunrise": 1601647503,
					"sunset": 1601689779,
					"type": 1
				},
				"timezone": -25200,
				"visibility": 4023,
				"weather": [
					{
						"description": "smoke",
						"icon": "50d",
						"id": 711,
						"main": "Smoke"
					},
					{
						"description": "haze",
						"icon": "50d",
						"id": 721,
						"main": "Haze"
					}
				],
				"wind": {
					"deg": 328,
					"speed": 1.42
				}
				}`),
			wantRes: &WeatherItem{
				Lat:             37.39,
				Lon:             -122.08,
				Description:     []string{"Smoke", "Haze"},
				CityName:        "Mountain View",
				ObservationTime: 1601662295,
				Temp:            28.87,
				MaxTemp:         31.67,
				MinTemp:         27,
				FeelsLike:       27.8,
				Humidity:        30,
			},
		},
		{
			name: "malformed response",
			lat:  1.0, lon: 2.0,
			cityName:      "Mountain View",
			apiRes:        []byte(`I am not a valid JSON`),
			apiStatusCode: http.StatusOK,
			wantErr:       true,
		},
		{
			name: "exceeded requests limit",
			lat:  1.0, lon: 2.0,
			cityName: "Mountain View",
			apiRes: []byte(`{
				"cod": 429,
				"message": "Your account is temporary blocked due to exceeding of requests limitation of your subscription type. 
				Please choose the proper subscription http://openweathermap.org/price"
			}`),
			apiStatusCode: http.StatusTooManyRequests,
			wantErr:       true,
		},
		{
			name: "invalid API key",
			lat:  1.0, lon: 2.0,
			cityName: "Mountain View",
			apiRes: []byte(`{
				"cod": 401,
				"message": "Invalid API key. Please see http://openweathermap.org/faq#error401 for more info."
			}`),
			apiStatusCode: http.StatusUnauthorized,
			wantErr:       true,
		},
		{
			name: "unreachable service",
			lat:  1.0, lon: 2.0,
			cityName:    "Mountain View",
			closeServer: true,
			apiRes:      []byte(``),
			wantErr:     true,
		},
		{
			name: "malformed URL",
			lat:  1.0, lon: 2.0,
			cityName:     "Mountain View",
			malformedURL: true,
			apiRes:       []byte(``),
			wantErr:      true,
		},
		{
			name: "city name is not valid",
			lat:  1.0, lon: 2.0,
			cityName: "Mountain View",
			apiRes: []byte(`{
				"cod": "404",
				"message": "city not found"
			}`),
			apiStatusCode: http.StatusNotFound,
			wantErr:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := newTestServer(test.apiStatusCode, test.apiRes)
			if test.closeServer {
				server.Close()
			} else {
				defer server.Close()
			}
			client := newTestAPIClient("apiKey", "metric", server, test.malformedURL)

			coordsRes, coordsErr := client.GetWeatherByCoords(test.lat, test.lon)
			compareResults(t, fmt.Sprintf("GetWeatherByCoords(%f, %f)", test.lat, test.lon), coordsRes, test.wantRes, coordsErr, test.wantErr)

			nameRes, nameErr := client.GetWeatherByCityName(test.cityName)
			compareResults(t, fmt.Sprintf("GetWeatherByCityName(%s)", test.cityName), nameRes, test.wantRes, nameErr, test.wantErr)
		})
	}
}
