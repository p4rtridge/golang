package main

import "fmt"

// omg it has so many ingredients
type Burger struct {
	buns      *string
	patty     *string
	cheese    *string
	vegetable *string
}

func (b *Burger) setBuns(bunStyle string) {
	b.buns = &bunStyle
}

func (b *Burger) setPatty(pattyStyle string) {
	b.patty = &pattyStyle
}

func (b *Burger) setCheese(cheeseStyle string) {
	b.cheese = &cheeseStyle
}

func (b *Burger) setVegetable(vegetableStyle string) {
	b.vegetable = &vegetableStyle
}

type BurgerBuilder struct {
	burger Burger
}

func NewBurgerBuilder() *BurgerBuilder {
	return &BurgerBuilder{
		burger: Burger{},
	}
}

func (bb *BurgerBuilder) addBuns(bunStyle string) *BurgerBuilder {
	bb.burger.setBuns(bunStyle)

	return bb
}

func (bb *BurgerBuilder) addPatty(pattyStyle string) *BurgerBuilder {
	bb.burger.setPatty(pattyStyle)

	return bb
}

func (bb *BurgerBuilder) addCheese(cheeseStyle string) *BurgerBuilder {
	bb.burger.setCheese(cheeseStyle)

	return bb
}

func (bb *BurgerBuilder) addVegetable(vegetableStyle string) *BurgerBuilder {
	bb.burger.setVegetable(vegetableStyle)

	return bb
}

func (bb BurgerBuilder) build() Burger {
	return bb.burger
}

func main() {
	burger := NewBurgerBuilder().addBuns("sesame").
		addCheese("swiss cheese").
		build()

	fmt.Println(burger)
}
