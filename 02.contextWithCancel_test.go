package golang_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// contoh goroutine leak (tidak berhenti)
func CreateCounterLeak() chan int {
	destination := make(chan int)
	go func() {
		defer close(destination)
		counter := 1
		// perulangan tak henti, dan mengirim data ke channel destination
		for {
			destination <- counter
			counter++
		}
	}()
	return destination
}

func TestContextWithCencelLeak(t *testing.T) {
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

	destination := CreateCounterLeak()

	// tidak ada kode untuk menghentikan goroutine
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	fmt.Println("Total Goroutine:", runtime.NumGoroutine()) // Goroutine leak (goroutine masih ada walau program sudah dihentikan)
}

// ############# ------------ #############

// contoh menghentikan goroutine dengan context
func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)
	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return // menggunakan return berhanti menghentikan for nya, kalau break cuman menghentikan select nya
			default:
				destination <- counter
				counter++
			}
		}
	}()
	return destination
}

func TestContextWith(t *testing.T) {
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	cancel()                    // mengirim sinyal cancel ke context
	time.Sleep(1 * time.Second) // untuk memastikan goroutine sudah berhenti
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())
}
