// Package greet 는 인사말을 만드는 작은 패키지입니다.
package greet

import "fmt"

// Hello 는 이름을 받아 인사말을 돌려줍니다.
// 이름이 대문자로 시작하므로 다른 패키지에서 greet.Hello 로 쓸 수 있습니다.
func Hello(name string) string {
	return decorate(fmt.Sprintf("%s님, 환영합니다", name))
}

// decorate 는 소문자로 시작하므로 greet 패키지 안에서만 쓸 수 있습니다.
func decorate(text string) string {
	return "** " + text + " **"
}
