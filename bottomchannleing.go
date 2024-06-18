package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// Goroutine untuk mengirim data ke ch1
	go func() {
		time.Sleep(1 * time.Second)
		ch1 <- 1
	}()

	// Goroutine untuk mengirim data ke ch2
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- 2
	}()

	for i := 0; i < 2; i++ {
		select {
		case value := <-ch1:
			fmt.Println("Received from ch1:", value)
		case value := <-ch2:
			fmt.Println("Received from ch2:", value)
		}
	}
	// Output:
	// Received from ch1: 1
	// Received from ch2: 2
}
