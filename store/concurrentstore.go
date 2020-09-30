package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/pablotrinidad/weatherreport/store/openweather"
)

// ConcurrentStore is a concurrent Store implementation with a built-in cache layer.
type ConcurrentStore struct {
	// results is a concurrent-safe map for storing API results. It acts as a cache as it prevents
	// HTTP requests from being duplicated.
	results *sync.Map

	// ow is an Open Weather API client.
	ow openweather.API
}

func NewConcurrentStore(ow openweather.API) Store {
	return &ConcurrentStore{ow: ow, results: &sync.Map{}}
}

// GetWeatherReport returns the weather report for the given airports on the current date and time.
// The returned map contains the airport code as the key and a weather report instance as value.
func (s ConcurrentStore) GetWeatherReport(airports []Airport) (map[string]WeatherReport, error) {
	if err := s.fetchConcurrently(airports); err != nil {
		return nil, err
	}
	data := make(map[string]WeatherReport)
	s.results.Range(func(key, value interface{}) bool {
		code := key.(string)
		report := value.(openweather.WeatherItem)
		data[code] = WeatherReport{
			Lat:             report.Lat,
			Lon:             report.Lon,
			Temp:            report.Temp,
			FeelsLike:       report.FeelsLike,
			Humidity:        report.Humidity,
			ObservationTime: time.Unix(int64(report.ObservationTime), 0),
		}
		return true
	})
	return data, nil
}

func (s ConcurrentStore) fetchConcurrently(airports []Airport) error {
	errors := make(chan error)
	wgDone := make(chan bool)

	var wg sync.WaitGroup
	for _, a := range airports {
		wg.Add(1)
		go func(a Airport) {
			defer wg.Done()
			if _, ok := s.results.Load(a.Code); ok {
				return
			}
			report, err := s.ow.GetCurrentWeather(a.Latitude, a.Longitude)
			if err != nil {
				errors <- err
			}
			s.results.Store(a.Code, report)
		}(a)
	}

	// Goroutine to wait until wait group is done.
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either wait group is done or an error is received trough the errors channel.
	select {
	case <-wgDone:
		// All goroutines finished successfully
		return nil
	case err := <-errors:
		close(errors)
		return fmt.Errorf("an error ocurred while communicating with Open Weather API: %v", err)
	}
}
