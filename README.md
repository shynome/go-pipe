## Usage

```go
package main

import (
	"io"
	"os"

	"github.com/shynome/go-pipe"
)

func main() {
	p := pipe.Line(
		pipe.HTTP("http://ip.sb/", nil),
		func(s pipe.State) error { io.Copy(os.Stdout, s.Stdin()); return nil },
	)
	pipe.Run(p)
}
```

## API

API inspired from [go-pipe](https://github.com/go-pipe/pipe) [article](https://labix.org/pipe)
