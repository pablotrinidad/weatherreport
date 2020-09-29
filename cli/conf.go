package main

import (
	"fmt"
	"os"
)

type Config struct {
	openweatherAPIKey string
}

func getConfig() (*Config, error) {
	owAPIKey, ok := os.LookupEnv("OPENWEATHER_API_KEY")
	if !ok {
		return nil, fmt.Errorf("failed to load OpenWeather API key from environment variables")
	}
	return &Config{openweatherAPIKey: owAPIKey}, nil
}
