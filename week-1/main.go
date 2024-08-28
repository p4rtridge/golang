package main

import (
	"fmt"
	"sync"
)

var (
	counter int
	mutex   sync.Mutex
)

func increment_with_mutex(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		// Lock the mutex before accessing the counter
		mutex.Lock()
		counter++
		// Unlock the mutex after the counter is incremented
		mutex.Unlock()
	}
}

func increment_without_mutex(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		counter++
	}
}

func prevent_race_condition() {
	var wg sync.WaitGroup

	// Spawn 10 goroutines to increment the counter
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go increment_with_mutex(&wg)
	}

	wg.Wait()
	fmt.Println("Final counter value:", counter)
}

func race_condition() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go increment_without_mutex(&wg)
	}

	wg.Wait()
	fmt.Println("Final counter value:", counter)
}

func main() {
	race_condition()
	prevent_race_condition()
}
