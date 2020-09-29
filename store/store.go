package store

// Store exposes a series of methods for querying weather information of specific cities.
type Store interface {
	// GetWeatherReport for a specific latitude/longitude pair at a given date and time.
	GetWeatherReport(lat, lon float64, datetime string) (WeatherReport, error)
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
	Humidity uint
	// ObservationTime of the report in ISO 8601.
	ObservationTime string
}
