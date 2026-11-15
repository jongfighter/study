package practice

// Clamp 는 value 를 lo 이상 hi 이하 범위로 맞춰 돌려줍니다.
func Clamp(value, lo, hi int) int {
	if value < lo {
		return lo
	}
	if value > hi {
		return hi
	}
	return value
}
