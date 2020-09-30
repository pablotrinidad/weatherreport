package openweather

import (
	"fmt"
)

// APIMockClient is a OpenWeather API mock implementation.
type APIMockClient struct {
	// GetCurrentWeatherCalls is the number of times GetCurrentWeather have been called.
	GetCurrentWeatherCalls uint

	// FailNext makes the next method call return an error if set to true, after failing,
	// value will be toggled back to false.
	FailNext bool

	weatherItem WeatherItem
}

func NewAPIMockClient(fixedWeatherItem WeatherItem) *APIMockClient {
	return &APIMockClient{
		GetCurrentWeatherCalls: 0,
		FailNext:               false,
		weatherItem:            fixedWeatherItem,
	}
}

// GetCurrentWeather returns an arbitrary weather item response and increments the calls registry.
func (c *APIMockClient) GetCurrentWeather(lat, lon float64) (*WeatherItem, error) {
	c.GetCurrentWeatherCalls++
	if c.FailNext {
		c.FailNext = false
		return nil, fmt.Errorf("expected fail after c.FailNext was set to true")
	}
	// Copy is needed since we don't want property modifications
	item := c.weatherItem
	return &item, nil
}
