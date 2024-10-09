package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	for {
		fmt.Println("hello to your mom")
		time.Sleep(10 * time.Second)
	}

	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, os.Interrupt, syscall.SIGTERM)
	<-cancel

	log.Println("[LOG] Shutting down...")
}
