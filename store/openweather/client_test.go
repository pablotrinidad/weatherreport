package openweather

import (
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

func TestAPIClient_GetCurrentWeather(t *testing.T) {
	tests := []struct {
		name          string
		lat, lon      float64
		apiRes        []byte
		apiStatusCode int
		malformedURL  bool
		wantRes       *WeatherItem
		wantErr       bool
		closeServer   bool
	}{
		{
			name: "successful call",
			lat:  37.39, lon: -122.08,
			apiRes: []byte(`{
				"coord": {
					"lat": 37.39,
					"lon": -122.08
				},
				"main": {
					"temp": 282.55,
					"feels_like": 281.86,
					"temp_min": 280.37,
					"temp_max": 284.26,
					"pressure": 1023,
					"humidity": 100
			  	},
				"dt": 1560350645,
				"timezone": -25200,
				"id": 420006353,
				"name": "Mountain View"
			}`),
			wantRes: &WeatherItem{
				Lat:             37.39,
				Lon:             -122.08,
				Temp:            282.55,
				FeelsLike:       281.86,
				Humidity:        100,
				ObservationTime: 1560350645,
			},
			apiStatusCode: http.StatusOK,
		},
		{
			name: "malformed response",
			lat:  1.0, lon: 2.0,
			apiRes:        []byte(`I am not a valid JSON`),
			apiStatusCode: http.StatusOK,
			wantErr:       true,
		},
		{
			name: "exceeded requests limit",
			lat:  1.0, lon: 2.0,
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
			closeServer: true,
			apiRes:      []byte(``),
			wantErr:     true,
		},
		{
			name: "malformed URL",
			lat:  1.0, lon: 2.0,
			malformedURL: true,
			apiRes:       []byte(``),
			wantErr:      true,
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
			gotRes, gotErr := client.GetCurrentWeather(test.lat, test.lon)
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("GetCurrentWeather(%f, %f) returned nil error, want error", test.lat, test.lon)
				}
				return
			}
			if gotErr != nil {
				t.Fatalf("GetCurrentWeather(%f, %f) returned unexpected error: %v", test.lat, test.lon, gotErr)
			}
			if diff := cmp.Diff(gotRes, test.wantRes); diff != "" {
				t.Errorf("GetCurrentWeather(%f, %f): %v, want %v\ngot -> want diff: %s", test.lat, test.lon, gotRes, test.wantRes, diff)
			}
		})
	}
}


func TestAPIClient_GetWeatherByCityName(t *testing.T) {
	tests := []struct {
		name          string
		cityName     string
		apiRes        []byte
		apiStatusCode int
		malformedURL  bool
		wantRes       *WeatherItem
		wantErr       bool
		closeServer   bool
	}{
		{
			name: "successful call",
			cityName: "Mountain View",
			apiRes: []byte(`{
				"coord": {
					"lat": 37.39,
					"lon": -122.08
				},
				"dt": 1601657590,
				"id": 5375480,
				"main": {
					"feels_like": 23.25,
					"humidity": 43,
					"pressure": 1016,
					"temp": 23.77,
					"temp_max": 25,
					"temp_min": 21
				},
				"name": "Mountain View"
			}`),
			wantRes: &WeatherItem{
				Lat:             37.39,
				Lon:             -122.08,
				Temp:            23.77,
				FeelsLike:       23.25,
				Humidity:        43,
				ObservationTime: 1601657590,
			},
			apiStatusCode: http.StatusOK,
		},
		{
			name: "malformed response",
			cityName: "Mountain View",
			apiRes:        []byte(`I am not a valid JSON`),
			apiStatusCode: http.StatusOK,
			wantErr:       true,
		},
		{
			name: "exceeded requests limit",
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
			cityName: "Mountain View",
			closeServer: true,
			apiRes:      []byte(``),
			wantErr:     true,
		},
		{
			name: "malformed URL",
			cityName: "Mountain View",
			malformedURL: true,
			apiRes:       []byte(``),
			wantErr:      true,
		},
		{
			name: "city name is not valid",
			cityName: "Mountain View",
			apiRes: []byte(`{
				"cod": "404",
				"message": "city not found"
			}`),
			apiStatusCode: http.StatusNotFound,
			wantErr:      true,
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
			gotRes, gotErr := client.GetWeatherByCityName(test.cityName)
			if test.wantErr {
				if gotErr == nil {
					t.Fatalf("GetWeatherByCityName(%s) returned nil error, want error", test.cityName)
				}
				return
			}
			if gotErr != nil {
				t.Fatalf("GetWeatherByCityName(%s) returned unexpected error: %v", test.cityName, gotErr)
			}
			if diff := cmp.Diff(gotRes, test.wantRes); diff != "" {
				t.Errorf("GetWeatherByCityName(%s): %v, want %v\ngot -> want diff: %s", test.cityName, gotRes, test.wantRes, diff)
			}
		})
	}
}
