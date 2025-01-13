package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2_add_change := make(chan int)

	// Goroutine untuk mengirim data ke ch1
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- 1
	}()

	// Goroutine untuk mengirim data ke ch2_add_change
	go func() {
		time.Sleep(2 * time.Second)
		ch2_add_change <- 2
	}()

	for i := 0; i < 2; i++ {
		fmt.Println("select")
		select {
		case value := <-ch1:
			fmt.Println("Received from ch1:", value)
		case value := <-ch2_add_change:
			fmt.Println("Received from ch2_add_change:", value)
		}
	}
	// Output:
	// Received from ch1: 1
	// Received from ch2_add_change: 2
}
