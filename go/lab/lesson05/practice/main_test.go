package main

import (
	"errors"
	"math"
	"testing"
)

func TestTotalArea(t *testing.T) {
	shapes := []Shape{Rectangle{Width: 3, Height: 4}, Circle{Radius: 1}}
	want := 12 + math.Pi
	if got := totalArea(shapes); math.Abs(got-want) > 1e-9 {
		t.Errorf("totalArea = %v; %v 이어야 합니다", got, want)
	}
}

func TestParsePort(t *testing.T) {
	if port, err := parsePort("8080"); err != nil || port != 8080 {
		t.Errorf("parsePort(8080) = %d, %v", port, err)
	}
	for _, text := range []string{"abc", "0", "70000", ""} {
		if _, err := parsePort(text); !errors.Is(err, ErrInvalidPort) {
			t.Errorf("parsePort(%q) 는 ErrInvalidPort 를 감싸야 합니다: %v", text, err)
		}
	}
}
