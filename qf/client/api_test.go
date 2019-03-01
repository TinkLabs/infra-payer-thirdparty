package client

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

const QfAppCode = "------key----"

func TestQfClient(t *testing.T) {
	var n sync.WaitGroup

	for i := 0; i < 10000; i++ {
		n.Add(1)
		go func() {
			start := time.Now()
			qc := &API{}
			qc.Init(QfAppCode, nil)
			fmt.Printf("%s, %+v\n%#v\n%#v\n%#v\n", time.Since(start), qc, qc.Charges, qc.Refunds, qc.Querys)

			n.Done()
		}()
	}

	n.Wait()
}
