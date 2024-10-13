package main

import (
	"kitchen/services/kitchen/composer"
	"log"
)

func main() {
	if err := composer.ServeHTTP("0.0.0.0:6969"); err != nil {
		log.Fatalln(err)
	}
}
