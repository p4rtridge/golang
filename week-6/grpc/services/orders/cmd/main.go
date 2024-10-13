package main

import (
	"kitchen/services/orders/composer"
	"log"
)

func main() {
	go func() {
		if err := composer.ServeRPC("0.0.0.0:50051"); err != nil {
			log.Fatalln(err)
		}
	}()

	if err := composer.ServeHTTP("0.0.0.0:8080"); err != nil {
		log.Fatalln(err)
	}
}
