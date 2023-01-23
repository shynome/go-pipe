package pipe

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
)

type Pipe func(s State) error

type ExitCode int
type Stdin interface {
	io.Reader
}
type Stdout interface {
	io.ReadWriter
}

type State interface {
	Context() context.Context
	Stdin() Stdin
	Stdout() Stdout
	Stderr() Stdout
	Exit(code ExitCode)
	ExitCode() ExitCode
}

type PipeState struct {
	ctx      context.Context
	stdin    Stdin
	stdout   Stdout
	stderr   Stdout
	exitCode ExitCode
}

var _ State = (*PipeState)(nil)

func NewPipeState(l State) (s *PipeState) {
	if l == nil {
		s = &PipeState{
			ctx:    context.Background(),
			stdin:  nil,
			stdout: &bytes.Buffer{},
		}
		return
	}
	s = &PipeState{
		ctx:    l.Context(),
		stdin:  l.Stdout(),
		stdout: &bytes.Buffer{},
		stderr: l.Stderr(),
	}
	return
}

func (p *PipeState) Exit(code ExitCode)       { p.exitCode = code }
func (p *PipeState) ExitCode() ExitCode       { return p.exitCode }
func (p *PipeState) Context() context.Context { return p.ctx }
func (p *PipeState) Stdin() Stdin             { return p.stdin }
func (p *PipeState) Stdout() Stdout           { return p.stdout }
func (p *PipeState) Stderr() Stdout           { return p.stderr }

type PipeStdout struct {
	bufrw *bufio.ReadWriter
}

var _ Stdout = (*PipeStdout)(nil)

func NewStdout() *PipeStdout {
	_buf := make([]byte, 512)
	buf := bytes.NewBuffer(_buf)
	bufrw := bufio.NewReadWriter(
		bufio.NewReader(buf),
		bufio.NewWriter(buf),
	)
	return &PipeStdout{
		bufrw: bufrw,
	}
}
func (out *PipeStdout) Read(p []byte) (n int, err error) {
	return out.Read(p)
}
func (out *PipeStdout) Write(p []byte) (n int, err error) {
	return out.Write(p)
}
func (stdout *PipeStdout) Close() error {
	return nil
}

var ErrExited = errors.New("pipe exited")

func Line(pipes ...Pipe) Pipe {
	return func(s State) error {
		for _, pipe := range pipes {

			var err error

			var errChan = make(chan error, 0)
			go func(s State) { errChan <- pipe(s) }(s)
			ctx := s.Context()
			select {
			case err = <-errChan:
			case <-ctx.Done():
				err = ctx.Err()
			}
			if err != nil {
				return err
			}

			if code := s.ExitCode(); code != 0 {
				return fmt.Errorf("%w. exit code is %d", ErrExited, code)
			}

			s = NewPipeState(s)
		}
		return nil
	}
}

func Run(pipe Pipe) error {
	s := NewPipeState(nil)
	return pipe(s)
}

type Input struct {
	Context context.Context
	Stdin   Stdin
	Stderr  Stdout
}

func RunWith(input Input, pipe Pipe) error {
	s := NewPipeState(nil)
	if input.Context != nil {
		s.ctx = input.Context
	}
	if input.Stdin != nil {
		s.stdin = input.Stdin
	}
	if input.Stderr != nil {
		s.stderr = input.Stderr
	}
	return pipe(s)
}
