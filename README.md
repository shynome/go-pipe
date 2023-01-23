## Usage

```go
func main(){
  p := pipe.Line(
    func(s pipe.State) error {}
  )
  err := pipe.Run(p)
}
```

## Todo

- [ ] pipe stdout

## API

API inspired from [go-pipe](https://github.com/go-pipe/pipe) [article](https://labix.org/pipe)
