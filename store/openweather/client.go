package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	baseURL            = "https://api.openweathermap.org/data/2.5/"
	currentWeatherPath = "weather"
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
	return &APIClient{apiKey: apiKey, units: units, apiURL: baseURL, client: &http.Client{}}, nil
}

// GetWeatherByCoords returns the current weather at the given location.
// It mirrors https://openweathermap.org/current.
func (c *APIClient) GetWeatherByCoords(lat, lon float64) (*WeatherItem, error) {
	res, err := c.makeHTTPCall(currentWeatherPath, map[string]string{
		"lat":   fmt.Sprintf("%f", lat),
		"lon":   fmt.Sprintf("%f", lon),
		"units": c.units,
	})
	if err != nil {
		return nil, err
	}
	return c.parseSuccessfulResponse(res.Body)
}

// GetWeatherByCityName returns the current weather at the given city name.
func (c *APIClient) GetWeatherByCityName(cityName string) (*WeatherItem, error) {
	res, err := c.makeHTTPCall(currentWeatherPath, map[string]string{
		"q":     cityName,
		"units": c.units,
	})
	if err != nil {
		return nil, err
	}
	return c.parseSuccessfulResponse(res.Body)
}

type currentWeatherResponse struct {
	ObservationTime int                      `json:"dt"`
	Coordinates     weatherResponseCoords    `json:"coord"`
	Weather         []weatherResponseWeather `json:"weather"`
	Data            *WeatherItem             `json:"main"`
	CityName        string                   `json:"name"`
}

type weatherResponseCoords struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type weatherResponseWeather struct {
	Description string `json:"main"`
}

func (c *APIClient) parseSuccessfulResponse(content io.ReadCloser) (*WeatherItem, error) {
	data := currentWeatherResponse{}
	decoder := json.NewDecoder(content)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed parsing API response: %v", err)
	}
	item := data.Data
	item.Lat = data.Coordinates.Lat
	item.Lon = data.Coordinates.Lon
	item.CityName = data.CityName
	item.ObservationTime = data.ObservationTime

	item.Description = make([]string, len(data.Weather))
	for i, d := range data.Weather {
		item.Description[i] = d.Description
	}
	return item, nil
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
		return nil, handleError(res)
	}
	return res, nil
}

type apiError struct {
	Code    string `json:"cod"`
	Message string `json:"message"`
}

func handleError(res *http.Response) error {
	apiError := apiError{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&apiError); err != nil {
		return fmt.Errorf("OpenWeather API returned unexpected error")
	}
	switch res.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("resource not found: %q", apiError.Message)
	case http.StatusTooManyRequests:
		return fmt.Errorf("exceeded requests limit")
	case http.StatusUnauthorized:
		return fmt.Errorf("invalid API key")
	}
	return fmt.Errorf("unexpected error")
}
