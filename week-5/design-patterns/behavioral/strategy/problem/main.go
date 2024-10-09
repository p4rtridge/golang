package main

import (
	"fmt"
)

type Car struct{}

type Bike struct{}

type Ship struct{}

type Manufactor struct {
	manufactureType string
}

// hard to maintain
// violate S and O in SOLID
func (m Manufactor) Manufacture() {
	switch m.manufactureType {
	case "car":
		fmt.Println("Manufacturing car")
	case "bike":
		fmt.Println("Manufacturing bike")
	case "ship":
		fmt.Println("Manufacturing ship")
	}
}

func main() {
	manufactor := Manufactor{
		manufactureType: "car",
	}

	manufactor.Manufacture()
}
