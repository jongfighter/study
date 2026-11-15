package main

import "fmt"

func main() {
	var stock map[string]int
	fmt.Println(stock["apple"])
	stock["apple"] = 3
}
