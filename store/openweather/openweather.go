// Package openweather exposes an interface and client for making HTTP requests
// to Open Weather's REST API (https://openweathermap.org/api).
package openweather

// API is a https://openweathermap.org/api API client.
type API interface {
	// GetWeatherByCoords returns the current weather at the given location.
	// It mirrors https://openweathermap.org/current.
	GetWeatherByCoords(lat, lon float64) (*WeatherItem, error)

	// GetWeatherByCityName returns the current weather at the given city name.
	GetWeatherByCityName(cityName string) (*WeatherItem, error)
}

// WeatherItem holds weather information for a given observation time.
type WeatherItem struct {
	// Latitude of the report location.
	Lat float64
	// Longitude of the report location.
	Lon float64
	// Description is a human readable set of weather descriptions
	Description []string
	// CityName is the city name registered in the API dataset for the weather observation.
	CityName string
	// ObservationTime in UNIX time UTC
	ObservationTime int
	// Temp is the temperature in celsius.
	Temp float64 `json:"temp"`
	// MaxTemp is the maximum expected temperature for the observation time.
	MaxTemp float64 `json:"temp_max"`
	// MinTemp is the maximum expected temperature for the observation time.
	MinTemp float64 `json:"temp_min"`
	// FeelsLike in celsius.
	FeelsLike float64 `json:"feels_like"`
	// Humidity percentage.
	Humidity int `json:"humidity"`
}
