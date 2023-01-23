package pipe_test

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/lainio/err2/assert"
	"github.com/shynome/go-pipe"
)

func TestPipe(t *testing.T) {
	var word = "hello world"
	p := pipe.Line(
		func(s pipe.State) error {
			go func() {
				defer s.Stdout().Close()
				io.WriteString(s.Stdout(), word)
			}()
			return nil
		},
		func(s pipe.State) error {
			b, err := io.ReadAll(s.Stdin())
			if err != nil {
				t.Error(err)
				return err
			}
			assert.Equal(string(b), word)
			return nil
		},
		func(s pipe.State) error { s.Exit(1); return nil },
		func(s pipe.State) error {
			io.WriteString(s.Stderr(), word)
			return nil
		},
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	input := pipe.Input{Context: ctx}
	err := pipe.RunWith(input, p)
	if err != nil && !errors.Is(err, pipe.ErrExited) {
		t.Error(err)
		return
	}
}
