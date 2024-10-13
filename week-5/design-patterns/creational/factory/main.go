package main

import (
	"fmt"
)

type Burger struct {
	ingredients []string
}

func (b Burger) print() {
	fmt.Println(b.ingredients)
}

type BurgerFactory struct{}

func (bf BurgerFactory) createCheeseBurger() Burger {
	ingredients := []string{"bun", "cheese", "beef-party"}

	return Burger{
		ingredients,
	}
}

func (bf BurgerFactory) createDeluxeCheeseBurger() Burger {
	ingredients := []string{"bun", "tomato", "lettuce", "cheese", "beef-party"}

	return Burger{
		ingredients,
	}
}

// careful, we can't control the detail. How burger are made
func (bf BurgerFactory) createVeganBurger() Burger {
	ingredients := []string{"bun", "special-sauce", "beef-party"}

	return Burger{
		ingredients,
	}
}

// or we can create like this
func createBurger(burgerType string) Burger {
	switch burgerType {
	case "cheese":
		ingredients := []string{"bun", "cheese", "beef-party"}

		return Burger{
			ingredients,
		}
	case "deluxe":
		ingredients := []string{"bun", "tomato", "lettuce", "cheese", "beef-party"}

		return Burger{
			ingredients,
		}
	case "vegan":
		ingredients := []string{"bun", "special-sauce", "beef-party"}

		return Burger{
			ingredients,
		}
	}

	return Burger{}
}

func main() {
	burgerFactory := new(BurgerFactory)

	burgerFactory.createCheeseBurger().print()
	burgerFactory.createDeluxeCheeseBurger().print()
	burgerFactory.createVeganBurger().print()
}
