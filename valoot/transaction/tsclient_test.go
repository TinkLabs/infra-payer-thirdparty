package transaction

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestValootBackend(t *testing.T) {
	var n sync.WaitGroup

	for i := 0; i < 100000; i++ {
		n.Add(1)
		go func() {
			start := time.Now()
			client := getC()
			fmt.Printf("%s, %v\n", time.Since(start), client)

			n.Done()
		}()
	}

	n.Wait()
}
