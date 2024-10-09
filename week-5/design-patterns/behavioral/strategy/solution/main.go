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

type Manufactor struct {
	vehicle ManufactureVehicle
}

func (m *Manufactor) SetVehicle(vehicle ManufactureVehicle) {
	m.vehicle = vehicle
}

func (m Manufactor) Manufacture() {
	m.vehicle.Manufacture()
}

func main() {
	manufactor := Manufactor{}
	manufactor.SetVehicle(Bike{})

	manufactor.Manufacture()
}
