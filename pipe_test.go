package pipe_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/lainio/err2/assert"
	"github.com/shynome/go-pipe"
)

func TestPipe(t *testing.T) {
	var word = "hello world"
	p := pipe.Line(
		func(s pipe.State) error {
			io.WriteString(s.Stdout(), word)
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
	var buf bytes.Buffer
	input := pipe.Input{Stderr: &buf}
	err := pipe.RunWith(input, p)
	assert.Equal(string(buf.Bytes()), "")
	if err != nil && !errors.Is(err, pipe.ErrExited) {
		t.Error(err)
		return
	}
}
