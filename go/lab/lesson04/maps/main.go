package main

import (
	"fmt"
	"maps"
	"slices"
)

func main() {
	stock := map[string]int{"apple": 3, "pear": 0}
	stock["kiwi"] = 5
	fmt.Println(stock, len(stock))

	count, ok := stock["banana"]
	fmt.Println("banana:", count, ok)

	if count, ok := stock["apple"]; ok {
		fmt.Println("apple 재고:", count)
	}

	delete(stock, "pear")

	for _, name := range slices.Sorted(maps.Keys(stock)) {
		fmt.Println(name, stock[name])
	}

	tags := map[string][]string{}
	tags["go"] = append(tags["go"], "backend")
	tags["go"] = append(tags["go"], "cli")
	fmt.Println(tags)

	var missing map[string]int
	fmt.Println("nil 맵 읽기:", missing["x"], len(missing))
}
