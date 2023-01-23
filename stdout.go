package pipe

import (
	"io"
)

type PipeStdout struct {
	r *io.PipeReader
	w *io.PipeWriter
}

var _ Stdout = (*PipeStdout)(nil)

func NewStdout() *PipeStdout {
	r, w := io.Pipe()
	return &PipeStdout{
		r: r, w: w,
	}
}
func (out *PipeStdout) Write(p []byte) (n int, err error) {
	n, err = out.w.Write(p)
	return
}
func (out *PipeStdout) Close() error {
	return out.w.Close()
}
func (out *PipeStdout) AsStdin() Stdin {
	return out.r
}
