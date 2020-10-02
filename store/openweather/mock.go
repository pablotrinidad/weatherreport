package openweather

import (
	"fmt"
)

// APIMockClient is a OpenWeather API mock implementation.
type APIMockClient struct {
	// APICalls stores the number of times all client methods have been called
	APICalls map[string]int

	// FailNext makes the next method call return an error if set to true, after failing,
	// value will be toggled back to false.
	FailNext bool

	weatherItem WeatherItem
}

func NewAPIMockClient(fixedWeatherItem WeatherItem) *APIMockClient {
	return &APIMockClient{
		APICalls:    map[string]int{},
		FailNext:    false,
		weatherItem: fixedWeatherItem,
	}
}

// GetWeatherByCoords returns an arbitrary weather item response and increments the calls registry.
func (c *APIMockClient) GetWeatherByCoords(_, _ float64) (*WeatherItem, error) {
	c.APICalls["coords"]++
	return c.produceResponse()
}

// GetWeatherByCityName returns an arbitrary weather item response and increments the calls registry.
func (c *APIMockClient) GetWeatherByCityName(_ string) (*WeatherItem, error) {
	c.APICalls["city"]++
	return c.produceResponse()
}

func (c *APIMockClient) produceResponse() (*WeatherItem, error) {
	if c.FailNext {
		c.FailNext = false
		return nil, fmt.Errorf("expected fail after c.FailNext was set to true")
	}
	// Copy is needed since we don't want any field modifications
	item := c.weatherItem
	return &item, nil
}
