package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("hello to your dad")
		time.Sleep(2 * time.Second)
	}
}
