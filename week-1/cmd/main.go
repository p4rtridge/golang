package main

import "github.com/partridge1307/week-1/pkg"

func main() {
	// Race condition
	pkg.RaceCondition()
	pkg.PreventRaceCondition()

	// Worker pool
	pkg.Workerpool()
}
