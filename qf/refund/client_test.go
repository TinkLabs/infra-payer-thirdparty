package refund

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/TinkLabs/payer-thirdparty/qf"
)

func TestQfRefundClient(t *testing.T) {
	qf.AppCode = "########key########"
	var n sync.WaitGroup

	for i := 0; i < 10000; i++ {
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
