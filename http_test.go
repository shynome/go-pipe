package pipe_test

import (
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	"github.com/shynome/go-pipe"
)

func TestHTTP(t *testing.T) {
	l := try.To1(net.Listen("tcp", "127.0.0.1:0"))
	defer l.Close()
	var word = "hello world"
	go http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(word))
	}))
	endpoint := "http://" + l.Addr().String()
	p := pipe.Line(
		pipe.HTTP(endpoint, nil),
		func(s pipe.State) error {
			b, err := io.ReadAll(s.Stdin())
			if err != nil {
				return err
			}
			assert.Equal(string(b), word)
			return nil
		},
	)
	err := pipe.Run(p)
	if err != nil {
		t.Error(err)
	}
}
