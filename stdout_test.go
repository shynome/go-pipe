package pipe

import (
	"bytes"
	"io"
	"testing"
)

func TestStdout(t *testing.T) {
	rw := NewStdout()
	var b = bytes.NewBuffer([]byte("555555555555"))
	ch := make(chan any)
	go func() {
		defer close(ch)
		s, err := io.ReadAll(rw.r)
		if err != nil {
			t.Error(err)
		}
		t.Log(s)
	}()
	go func() {
		defer rw.Close()
		_, err := io.Copy(rw, b)
		if err != nil {
			t.Error(err)
		}
	}()
	<-ch
}
