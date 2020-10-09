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
		return nil, fmt.Errorf("missing OpenWeather API key, please set envar OPENWEATHER_API_KEY to continue")
	}
	return &Config{openweatherAPIKey: owAPIKey}, nil
}
