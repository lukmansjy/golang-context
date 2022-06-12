package golang_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func Increment2(ctx context.Context) chan int {
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

func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5*time.Second))
	defer cancel()

	destination := Increment2(ctx)

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(1 * time.Second)

	fmt.Println("Total Goroutine:", runtime.NumGoroutine())
}
