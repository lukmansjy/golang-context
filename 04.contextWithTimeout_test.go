package golang_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func Increment(ctx context.Context) chan int {
	destination := make(chan int)
	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) // simulasi slow
			}
		}
	}()
	return destination
}

func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 3*time.Second)

	/*
		cancel tetap di panggil walau sudah pakai timeout. istilahnya time out itu dibatalkan ketika melebih waktu.
		dan kalau lebih prosesnya lebih cepat tetap harus jalankan cancel.
		untuk memastikan goroutine background yang akan jado leak
	*/
	defer cancel()

	destination := Increment(ctx)

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(1 * time.Second)

	fmt.Println("Total Goroutine:", runtime.NumGoroutine())
}
