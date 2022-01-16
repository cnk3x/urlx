package form

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

type (
	Process = func(resp *http.Response, body io.ReadCloser) error // 响应处理器
	Body    = func() (contentType string, body io.Reader, err error)
)

// Decode 处理JSON响应
func Decode(out any) Process {
	return func(resp *http.Response, body io.ReadCloser) error {
		defer body.Close()
		return json.NewDecoder(body).Decode(out)
	}
}

// SendJSON 提交JSON
func Encode(in any) Body {
	return func() (contentType string, body io.Reader, err error) {
		contentType = "application/x-www-form-urlencoded; charset=utf-8"
		switch o := in.(type) {
		case io.Reader:
			body = o
		case []byte:
			body = bytes.NewReader(o)
		case string:
			body = strings.NewReader(o)
		case bytes.Buffer:
			body = &o
		case *bytes.Buffer:
			body = o
		case url.Values:
			body = strings.NewReader(o.Encode())
		case *url.Values:
			body = strings.NewReader(o.Encode())
		case map[string]string:
			values := url.Values{}
			for k, v := range o {
				values.Set(k, v)
			}
			body = strings.NewReader(values.Encode())
		default:
			if r, ok := o.(io.Reader); ok {
				body = r
			} else {
				var values url.Values
				if values, err = query.Values(in); err == nil {
					body = strings.NewReader(values.Encode())
				}
			}
		}
		return
	}
}
