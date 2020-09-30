package main

import (
	"fmt"
	"os"
)

type DescentBeing interface {
	SayHello()
	Multiply(int, int) int
}

type Human struct {
	Name string
	Age  uint8
}

func (h *Human) SayHello() {
	fmt.Printf("hello, my names is %s", h.Name)
}

func (h *Human) Multiply(a, b int) int {
	return a * b
}

func possiblyUnsafeOperation(data ...int) (int, error) {
	fmt.Printf("%T, %v\n", data, data)
	var x DescentBeing
	x = &Human{}
	var f func(int, int) int
	f = func(a, b int) int { return a * b }
	f = func(_, _ int) int { return 0 }
	x.SayHello()
	fmt.Println(f(5, 3))
	return 0, nil
}

func main() {
	possiblyUnsafeOperation(1, 2, 3, 4, 5, 6, 7, 8, 8, 4)
	config, err := getConfig()
	if err != nil {
		fmt.Fprint(os.Stderr, fmt.Sprintf("Cannot start application: %v", err))
		os.Exit(1)
	}
	fmt.Println(config.openweatherAPIKey)
}
