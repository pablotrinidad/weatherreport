package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	baseURL            = "https://api.openweathermap.org/data/2.5"
	currentWeatherPath = "weather"
	oneCallPath        = "onecall"
)

// APIClient is an API implementation.
type APIClient struct {
	apiKey string
	apiURL string
	units  string
	client *http.Client
}

// NewAPIClient returns an Open Weather API client that uses the given API key and units system.
func NewAPIClient(apiKey, units string) (*APIClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("got empty API key")
	}
	if _, ok := map[string]bool{"standard": true, "metric": true, "imperial": true}[units]; !ok {
		return nil, fmt.Errorf("got invalid units value %s, want one of standard, metric, or imperial", units)
	}
	return &APIClient{apiKey: apiKey, units: units, apiURL: baseURL}, nil
}

type currentWeatherResponse struct {
	ObservationTime int          `json:"dt"`
	Data            *WeatherItem `json:"main"`
}

// GetCurrentWeather returns the current weather at the given location.
// It mirrors https://openweathermap.org/current.
func (c *APIClient) GetCurrentWeather(lat, lon float64) (*WeatherItem, error) {
	latS, lonS := fmt.Sprintf("%f", lat), fmt.Sprintf("%f", lon)
	res, err := c.makeHTTPCall(currentWeatherPath, map[string]string{"lat": latS, "lon": lonS})
	if err != nil {
		return nil, fmt.Errorf("failed making HTTP call: %v", err)
	}

	data := currentWeatherResponse{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed parsing API response: %v", err)
	}

	item := data.Data
	item.ObservationTime = data.ObservationTime
	return item, nil
}

func (c *APIClient) OneCall(lat, lon float64) (*OneCallResponse, error) {
	panic("implement me")
}

// makeHTTPCall performs an HTTP GET request to Open Weather's REST API using API access token.
func (c *APIClient) makeHTTPCall(path string, q map[string]string) (*http.Response, error) {
	base, err := url.Parse(c.apiURL)
	if err != nil {
		return nil, err
	}
	base.Path += path
	params := url.Values{}
	q["appid"] = c.apiKey
	for k, v := range q {
		params.Add(k, v)
	}
	base.RawQuery = params.Encode()

	res, err := c.client.Get(base.String())
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-succesful status code %d", res.StatusCode)
	}
	return res, nil
}
