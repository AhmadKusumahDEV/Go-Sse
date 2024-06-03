package main

import (
	"fmt"
	"time"
)

func main() {
	// Membuat channel
	done := make(chan bool)

	// Goroutine untuk melakukan tugas dan kemudian mengirim sinyal selesai
	go func() {
		fmt.Println("Goroutine memulai tugas")
		time.Sleep(2 * time.Second) // Simulasi tugas
		fmt.Println("Goroutine selesai tugas")
		done <- true
	}()

	// Menunggu sinyal selesai dari goroutine
	fmt.Println("Main function menunggu")
	let := <-done
	fmt.Println("Goroutine selesai", let)
	fmt.Println("Main function selesai menunggu")
}
