// Package exercises는 레슨 07 연습 문제의 정답 예시입니다.
package exercises

// ParallelSum은 nums를 parts개 조각으로 나눠 고루틴마다 합을 구한 뒤 모두 더합니다.
func ParallelSum(nums []int, parts int) int {
	if len(nums) == 0 {
		return 0
	}
	if parts < 1 {
		parts = 1
	}
	if parts > len(nums) {
		parts = len(nums)
	}

	size := (len(nums) + parts - 1) / parts
	sums := make(chan int)
	started := 0
	for start := 0; start < len(nums); start += size {
		end := min(start+size, len(nums))
		started++
		go func(chunk []int) {
			total := 0
			for _, n := range chunk {
				total += n
			}
			sums <- total
		}(nums[start:end])
	}

	total := 0
	for range started {
		total += <-sums
	}
	return total
}
