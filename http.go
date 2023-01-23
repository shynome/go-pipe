package pipe

import (
	"fmt"
	"io"
	"net/http"
)

var client *http.Client = http.DefaultClient

func SetHTTPCleint(nclient *http.Client) { client = nclient }

func HTTP(path string, body io.Reader) Pipe {
	req, err := http.NewRequest(http.MethodGet, path, body)
	req.Header.Set("User-Agent", "curl-go-pipe/0.0.1")
	return func(s State) error {
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if code := resp.StatusCode; !(200 <= code && code < 300) {
			return fmt.Errorf("statu code is not 2xx, got %d", code)
		}
		go func() {
			defer s.Stdout().Close()
			io.Copy(s.Stdout(), resp.Body)
		}()
		return nil
	}
}
