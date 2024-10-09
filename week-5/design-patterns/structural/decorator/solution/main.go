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

type ManufactorDecoratorWithHeap struct {
	core       *ManufactorDecoratorWithHeap
	manufactor ManufactureVehicle
}

func NewManufactor(vehicle ManufactureVehicle) ManufactorDecoratorWithHeap {
	return ManufactorDecoratorWithHeap{
		core:       nil,
		manufactor: vehicle,
	}
}

func (m ManufactorDecoratorWithHeap) Add(vehicle ManufactureVehicle) ManufactorDecoratorWithHeap {
	return ManufactorDecoratorWithHeap{
		core:       &m,
		manufactor: vehicle,
	}
}

func (m ManufactorDecoratorWithHeap) Manufacture() {
	m.manufactor.Manufacture()

	if m.core != nil {
		m.core.Manufacture()
	}
}

type Service struct {
	s ManufactureVehicle
}

func (s Service) Manufacture() {
	s.s.Manufacture()
}

func main() {
	m := NewManufactor(Car{}).Add(Bike{}).Add(Ship{})

	s := Service{
		s: m,
	}

	s.Manufacture()
}
