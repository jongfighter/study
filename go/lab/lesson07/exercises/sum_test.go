package exercises

import "testing"

func TestParallelSum(t *testing.T) {
	thousand := make([]int, 1000)
	for i := range thousand {
		thousand[i] = i + 1
	}

	tests := []struct {
		name  string
		nums  []int
		parts int
		want  int
	}{
		{"빈 슬라이스", nil, 4, 0},
		{"조각 하나", []int{1, 2, 3}, 1, 6},
		{"조각이 원소보다 많음", []int{1, 2, 3}, 10, 6},
		{"조각 수 0은 1로 취급", []int{4, 5}, 0, 9},
		{"1부터 1000까지", thousand, 7, 500500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParallelSum(tt.nums, tt.parts); got != tt.want {
				t.Errorf("ParallelSum(%d개, %d조각) = %d, want %d", len(tt.nums), tt.parts, got, tt.want)
			}
		})
	}
}
