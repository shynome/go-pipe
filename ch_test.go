package pipe_test

import (
	"testing"
	"time"

	"github.com/shynome/go-pipe"
)

func TestCh(t *testing.T) {
	exec := pipe.Go(
		func(in <-chan any, out chan<- any) {
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second)
				out <- i
			}
		},
		func(in <-chan any, out chan<- any) {
			for v := range in {
				out <- v.(int) * 2
			}
		},
	)
	var out = make(chan any)
	exec(nil, out)
	for v := range out {
		t.Log(v)
	}
}
