package store

import "time"

// Store exposes a series of methods for querying weather information of specific cities.
// It abstracts away cache layer and API access.
type Store interface {
	// GetWeatherReport returns the weather report for the given airports on the current date and time.
	// The returned map contains the airport code as the key and a weather report instance as value.
	GetWeatherReport([]Airport) (map[string]WeatherReport, error)
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
	// Temperature in celsius with two decimals precision.
	Temp float64
	// FeelsLike is the weather perceived temp in celsius with two decimals precision.
	FeelsLike float64
	// Humidity percentage.
	Humidity int
	// ObservationTime of the report in ISO 8601.
	ObservationTime time.Time
}
