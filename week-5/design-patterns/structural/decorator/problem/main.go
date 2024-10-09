package main

import "fmt"

type ManufactureVehicle interface {
	Manufacture()
}

type Car struct{}

func (c Car) Manufacture() {
	fmt.Println("Manufacturing car")
}

type Bike struct{}

func (b Bike) Manufacture() {
	fmt.Println("Manufacturing bike")
}

type Ship struct{}

func (s Ship) Manufacture() {
	fmt.Println("Manufacturing ship")
}

// duplicate
// wtf is this
type ManufactorBikeCar struct {
	bike Bike
	car  Car
}

func (m ManufactorBikeCar) Manufacture() {
	m.bike.Manufacture()
	m.car.Manufacture()
}

type ManufactorCarBike struct {
	bike Bike
	car  Car
}

func (m ManufactorCarBike) Manufacture() {
	m.car.Manufacture()
	m.bike.Manufacture()
}

func main() {
}
