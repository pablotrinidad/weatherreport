// Package openweather exposes an interface and client for making HTTP requests
// to Open Weather's REST API (https://openweathermap.org/api).
package openweather

// API is a https://openweathermap.org/api API client.
type API interface {
	// GetCurrentWeather returns the current weather at the given location.
	// It mirrors https://openweathermap.org/current.
	GetCurrentWeather(lat, lon float64) (*WeatherItem, error)

	// OneCall returns current, hourly and daily weather report at the given location.
	// It mirrors https://openweathermap.org/api/one-call-api with exclude: [minutely, alerts].
	OneCall(lat, lon float64) (*OneCallResponse, error)
}

// WeatherItem holds weather information for a given observation time.
type WeatherItem struct {
	// Temperature in celsius.
	Temp float64 `json:"temp"`
	// FeelsLike in celsius.
	FeelsLike float64 `json:"feels_like"`
	// Humidity percentage.
	Humidity int `json:"humidity"`
	// ObservationTime in UNIX time UTC
	ObservationTime int `json:"dt"`
}

// OneCallResponse represents the data returned by OpenWeather through their One Call endpoint.
type OneCallResponse struct {
	// Latitude of the queried location.
	Lat float64
	// Longitude of the queried location.
	Lon float64

	// Current weather at queried location.
	Current WeatherItem
	// Hourly weather items at queried location from 00:00 to 23:59 (UTC) on the same day of the request.
	Hourly []WeatherItem
	// Daily weather items at queries location for the next 7 days (plus current).
	Daily []WeatherItem
}
