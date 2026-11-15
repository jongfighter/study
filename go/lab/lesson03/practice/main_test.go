package main

import "testing"

func TestAverage(t *testing.T) {
	got, err := average(1, 2, 3, 4)
	if err != nil || got != 2.5 {
		t.Errorf("average(1, 2, 3, 4) = %v, %v; 2.5, nil 이어야 합니다", got, err)
	}
	if _, err := average(); err == nil {
		t.Error("값이 없으면 오류를 돌려줘야 합니다")
	}
}

func TestFizzBuzz(t *testing.T) {
	tests := map[int]string{1: "1", 3: "Fizz", 5: "Buzz", 15: "FizzBuzz", 14: "14"}
	for n, want := range tests {
		if got := fizzBuzz(n); got != want {
			t.Errorf("fizzBuzz(%d) = %q; %q 이어야 합니다", n, got, want)
		}
	}
}
