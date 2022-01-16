package json

import (
	"io"
	"net/http"

	"github.com/goccy/go-yaml"
)

type Process = func(resp *http.Response, body io.ReadCloser) error // 响应处理器

// Decode 处理yaml响应
func Decode(out any) Process {
	return func(resp *http.Response, body io.ReadCloser) error {
		defer body.Close()
		return yaml.NewDecoder(body).Decode(out)
	}
}
