package main

import "fmt"

type Animal interface {
	Speak()
}

type Dog struct{}

func (d Dog) Speak() {
	fmt.Println("Woan")
}

type Cat struct{}

func (c Cat) Speak() {
	fmt.Println("Meow")
}

type Fish struct{}

func (f Fish) Speak() {
	panic("I never heard that fish can speak")
}

type Bird struct{}

func (b Bird) Speak() {
	panic("I never heard that bird can speak")
}

type AnimalSpeak struct {
	animal Animal
}

func (a AnimalSpeak) Speak() {
	a.animal.Speak()
}

func main() {
	// hard to maintain, extension
	// can cause a bug
	animalType := "cat"

	var animal Animal
	switch animalType {
	case "dog":
		animal = new(Dog)
	case "cat":
		animal = new(Cat)
	case "fish":
		animal = new(Fish)
	case "bird":
		animal = new(Bird)
	}

	as := AnimalSpeak{
		animal,
	}

	as.Speak()
}
