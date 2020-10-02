package store

import (
	"time"

	"github.com/pablotrinidad/weatherreport/store/openweather"
)

const maxConcurrentRequestsPerMinute = 60

// ConcurrentStore is a concurrent Store implementation.
type ConcurrentStore struct {
	// ow is an Open Weather API client.
	ow openweather.API
}

func NewConcurrentStore(ow openweather.API) Store {
	return &ConcurrentStore{ow: ow}
}

// GetWeatherByAirportCode returns the weather report for the given airports on the current date and time.
// The returned map contains the airport code as the key and a weather report instance as value.
func (s *ConcurrentStore) GetWeatherByAirportCode(airports []Airport) map[string]WeatherReport {
	requests := make(map[string]func() (*openweather.WeatherItem, error))
	for i := range airports {
		a := airports[i]
		// Using a map (hash table) avoids repeating API requests for the same airport code.
		requests[a.Code] = func() (*openweather.WeatherItem, error) {
			return s.ow.GetWeatherByCoords(a.Latitude, a.Longitude)
		}
	}
	results := s.fetchConcurrently(requests)
	return s.parseResults(results)
}

// GetWeatherByCityName returns the weather report for each city name. The returned map contains
// the city name as key and a weather report instance as value.
func (s *ConcurrentStore) GetWeatherByCityName(cities []string) map[string]WeatherReport {
	requests := make(map[string]func() (*openweather.WeatherItem, error))
	for i := range cities {
		cityName := cities[i]
		// Using a map (hash table) avoids repeating API requests for the same city name.
		requests[cityName] = func() (*openweather.WeatherItem, error) {
			return s.ow.GetWeatherByCityName(cityName)
		}
	}
	results := s.fetchConcurrently(requests)
	return s.parseResults(results)
}

// parseResults .....
func (s *ConcurrentStore) parseResults(results map[string]*requestResult) map[string]WeatherReport {
	data := make(map[string]WeatherReport)
	for key, val := range results {
		if val.err != nil {
			data[key] = WeatherReport{
				Failed:      true,
				FailMessage: "failed getting weather data from external services",
			}
			continue
		}
		data[key] = WeatherReport{
			Lat:             val.data.Lat,
			Lon:             val.data.Lon,
			Description:     val.data.Description,
			CityName:        val.data.CityName,
			Temp:            val.data.Temp,
			MaxTemp:         val.data.MaxTemp,
			MinTemp:         val.data.MinTemp,
			FeelsLike:       val.data.FeelsLike,
			Humidity:        val.data.Humidity,
			ObservationTime: time.Unix(int64(val.data.ObservationTime), 0),
			Failed:          false,
		}
	}
	return data
}

type requestResult struct {
	data *openweather.WeatherItem
	key  string
	err  error
}

// fetchConcurrently something
func (s ConcurrentStore) fetchConcurrently(requests map[string]func() (*openweather.WeatherItem, error)) map[string]*requestResult {
	cn := make(chan *requestResult, len(requests))
	fns := make([]func(), len(requests))
	i := 0
	for key := range requests {
		f := requests[key]
		fns[i] = func() {
			report, err := f()
			cn <- &requestResult{data: report, err: err, key: key}
		}
		i++
	}

	// Concurrently process requests in batched of up to 60 requests per minute
	start := 0
	breakNext := false
	for !breakNext {
		end := start + maxConcurrentRequestsPerMinute
		if end > len(requests) {
			breakNext = true
			end = len(requests)
		}
		callConcurrent(fns[start:end])
		start = end
		if end-start == maxConcurrentRequestsPerMinute {
			time.Sleep(1 * time.Minute)
		}
	}
	close(cn)

	// Read results
	results := make(map[string]*requestResult, len(requests))
	for r := range cn {
		if _, ok := results[r.key]; !ok {
			results[r.key] = r
		}
	}
	return results
}
