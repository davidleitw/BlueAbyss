package main

import (
	"fmt"
)

func main() {
	for i := 0; i < 10; i++ {
		c := make(chan int, 2147483640)
		fmt.Println(c)
	}
}
