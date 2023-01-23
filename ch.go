package pipe

type GoPipe func(in <-chan any, out chan<- any)

func Go(pipes ...GoPipe) GoPipe {
	return func(in <-chan any, out chan<- any) {
		var po chan any
		for _, pipe := range pipes {
			po = make(chan any)
			go fpipe(pipe, in, po)
			in = po
		}
		go func() {
			if out == nil {
				for range po {
				}
				return
			}
			defer close(out)
			for v := range po {
				out <- v
			}
		}()
		return
	}
}

func fpipe(pipe GoPipe, in <-chan any, out chan<- any) {
	defer close(out)
	pipe(in, out)
}
