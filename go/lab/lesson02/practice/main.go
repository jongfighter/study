package main

import (
	"fmt"
	"os"

	"golab/lesson02/greet"
)

func main() {
	name := "개발자"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	fmt.Println(greet.Hello(name))
	fmt.Println(greet.Farewell(name))
}
