package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

var ErrInvalidPort = errors.New("잘못된 포트 번호")

func parsePort(text string) (int, error) {
	port, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("포트 %q: %w", text, ErrInvalidPort)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("포트 %d: %w", port, ErrInvalidPort)
	}
	return port, nil
}

func main() {
	shapes := []Shape{Rectangle{Width: 3, Height: 4}, Circle{Radius: 1}}
	fmt.Printf("도형 넓이 합: %.2f\n", totalArea(shapes))

	for _, text := range []string{"8080", "abc", "70000"} {
		port, err := parsePort(text)
		fmt.Println(port, err, errors.Is(err, ErrInvalidPort))
	}
}
