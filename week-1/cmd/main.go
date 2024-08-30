package main

import (
	"fmt"

	"github.com/partridge1307/week-1/pkg"
)

func main() {
	// Race condition
	fmt.Println("Race condition: ")
	pkg.RaceCondition()
	pkg.PreventRaceCondition()

	// Worker pool
	fmt.Println("\n\nWorker pool")
	pkg.Workerpool()
}
