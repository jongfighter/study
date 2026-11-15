package practice

import (
	"fmt"
	"testing"
)

func TestClamp(t *testing.T) {
	tests := []struct {
		name  string
		value int
		lo    int
		hi    int
		want  int
	}{
		{"범위 안", 5, 0, 10, 5},
		{"최솟값 아래", -3, 0, 10, 0},
		{"최댓값 위", 42, 0, 10, 10},
		{"경계값 lo", 0, 0, 10, 0},
		{"경계값 hi", 10, 0, 10, 10},
		{"음수 범위", -5, -10, -1, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp(tt.value, tt.lo, tt.hi); got != tt.want {
				t.Errorf("Clamp(%d, %d, %d) = %d; %d 이어야 합니다", tt.value, tt.lo, tt.hi, got, tt.want)
			}
		})
	}
}

func ExampleClamp() {
	fmt.Println(Clamp(42, 0, 10))
	// Output: 10
}

func BenchmarkClamp(b *testing.B) {
	for b.Loop() {
		Clamp(42, 0, 10)
	}
}
