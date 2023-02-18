package multiball

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

const (
	workers = 10
	tasks = 100
)

func TestMultiball(t *testing.T) {
	multiball := NewMultiball(0, func(next int64) {
		fmt.Printf("[%v] next <- %v\n", time.Now(), next)
	})
	defer multiball.Close()

	channel := make(chan int64)

	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func(worker int) {
			for task := range channel {
				time.Sleep(time.Duration(rand.Float64() * float64(time.Second)))
				fmt.Printf("[%v] finished %v\n", time.Now(), task)
				multiball.Finish <- task
			}
			wg.Done()
		}(worker)
	}

	for task := 0; task < tasks; task++ {
		channel <- int64(task)
	}
	close(channel)

	wg.Wait()
}
