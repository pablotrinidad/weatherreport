package openweather

import "C"
import (
	"fmt"
)

// APIMockClient is a OpenWeather API mock implementation.
type APIMockClient struct {
	// FailNext makes the next method call return an error if set to true.
	FailNext bool

	weatherItem WeatherItem
}

func NewAPIMockClient(fixedWeatherItem WeatherItem) *APIMockClient {
	return &APIMockClient{
		FailNext:    false,
		weatherItem: fixedWeatherItem,
	}
}

// GetWeatherByCoords returns an arbitrary weather item response.
func (c *APIMockClient) GetWeatherByCoords(_, _ float64) (*WeatherItem, error) {
	return c.produceResponse()
}

// GetWeatherByCityName returns an arbitrary weather item response.
func (c *APIMockClient) GetWeatherByCityName(_ string) (*WeatherItem, error) {
	return c.produceResponse()
}

func (c *APIMockClient) produceResponse() (*WeatherItem, error) {
	if c.FailNext {
		return nil, fmt.Errorf("expected fail after c.FailNext was set to true")
	}
	// Copy is needed since we don't want any field modifications
	item := c.weatherItem
	return &item, nil
}
