package main

import (
	"fmt"
	"os"
)

func main() {
	config, err := getConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("Cannot start application: %v", err))
		os.Exit(1)
	}
	fmt.Println(config.openweatherAPIKey)
}
