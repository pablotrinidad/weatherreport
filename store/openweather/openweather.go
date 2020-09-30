// Package openweather exposes an interface and client for making HTTP requests
// to Open Weather's REST API (https://openweathermap.org/api).
package openweather

// API is a https://openweathermap.org/api API client.
type API interface {
	// GetCurrentWeather returns the current weather at the given location.
	// It mirrors https://openweathermap.org/current.
	GetCurrentWeather(lat, lon float64) (*WeatherItem, error)
}

// WeatherItem holds weather information for a given observation time.
type WeatherItem struct {
	// Latitude of the report location.
	Lat float64
	// Longitude of the report location.
	Lon float64
	// Temperature in celsius.
	Temp float64 `json:"temp"`
	// FeelsLike in celsius.
	FeelsLike float64 `json:"feels_like"`
	// Humidity percentage.
	Humidity int `json:"humidity"`
	// ObservationTime in UNIX time UTC
	ObservationTime int `json:"dt"`
}
