package store

import "time"

// Store exposes a series of methods for querying weather information of specific cities.
// It abstracts away cache layer and API access.
type Store interface {
	// GetWeatherReport returns the weather report for the given airports on the current date and time.
	// The returned map contains the airport code as the key and a weather report instance as value.
	GetWeatherByAirportCode([]Airport) map[string]WeatherReport

	// GetWeatherByCityName returns the weather report for each city name. The returned map contains
	// the city name as key and a weather report instance as value.
	GetWeatherByCityName([]string) map[string]WeatherReport

	// GetAPIUsage returns OpenWeather API usage statistics.
	GetAPIUsage() APIUsage
}

// Airport data.
type Airport struct {
	// Airport code, e.g: code for Mexico City International Airport is MEX.
	Code string
	// Latitude of the airport location.
	Latitude float64
	// Longitude of the airport location.
	Longitude float64
}

// WeatherReport holds the information of an weather query for a specific latitude, longitude pair.
type WeatherReport struct {
	// Latitude of the report location.
	Lat float64
	// Longitude of the report location.
	Lon float64
	// Description is a human readable set of weather descriptions
	Description []string
	// CityName is the city name registered in the API dataset for the weather observation.
	CityName string
	// Temperature in celsius with two decimals precision.
	Temp float64 `json:"temp"`
	// MaxTemp is the maximum expected temperature for the observation time.
	MaxTemp float64 `json:"temp_max"`
	// MinTemp is the maximum expected temperature for the observation time.
	MinTemp float64 `json:"temp_min"`
	// FeelsLike in celsius.
	FeelsLike float64 `json:"feels_like"`
	// Humidity percentage.
	Humidity        int `json:"humidity"`
	ObservationTime time.Time
	// Failed indicates that the API request was unsuccessful
	Failed bool
	// FailMessage is the reason of failure.
	FailMessage string
}

// APIUsage contains usage statistics.
type APIUsage struct {
	// SuccessfulCalls count.
	SuccessfulCalls uint
	// FailedCalls count.
	FailedCalls uint
}
