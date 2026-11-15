package greet

import "fmt"

// Farewell 는 작별 인사를 돌려줍니다.
func Farewell(name string) string {
	return decorate(fmt.Sprintf("%s님, 다음에 또 만나요", name))
}
