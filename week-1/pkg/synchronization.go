package pkg

import (
	"fmt"
	"sync"
)

func RaceCondition() {
	var (
		wg      sync.WaitGroup
		counter = 0
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			counter++
		}()
	}

	wg.Wait()

	fmt.Printf("Counter without Mutex: %d\n", counter)
}

func PreventRaceCondition() {
	var (
		wg      sync.WaitGroup
		counter = 0
		mu      sync.Mutex
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			mu.Lock()

			defer wg.Done()
			defer mu.Unlock()

			counter++
		}()
	}

	wg.Wait()

	fmt.Printf("Counter with Mutex: %d\n", counter)
}
