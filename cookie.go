package urlx

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	HeaderRequestCookie  = "Cookie"     // Request Cookie
	HeaderResponseCookie = "Set-Cookie" // Response Cookie
)

// ReadCookie 从响应读取Cookie
func ReadCookie(read func(cookies []*http.Cookie) error) ProcessMw {
	return func(next Process) Process {
		return func(resp *http.Response) error {
			if err := read(resp.Cookies()); err != nil {
				return err
			}
			return next(resp)
		}
	}
}

// CookieAddString 添加Cookie到请求
func CookieAddString(cookies ...string) HeaderOption {
	return func(headers http.Header) {
		for _, s := range cookies {
			if s != "" {
				if c := headers.Get(HeaderRequestCookie); c != "" {
					headers.Set(HeaderRequestCookie, c+"; "+s)
				} else {
					headers.Set(HeaderRequestCookie, s)
				}
			}
		}
	}
}

// CookieAdd 添加Cookie到请求
func CookieAdd(cookies ...*http.Cookie) HeaderOption {
	return func(headers http.Header) {
		for _, cookie := range cookies {
			if cookie != nil {
				s := fmt.Sprintf("%s=%s", sanitizeCookieName(cookie.Name), sanitizeCookieValue(cookie.Value))
				if c := headers.Get(HeaderRequestCookie); c != "" {
					headers.Set(HeaderRequestCookie, c+"; "+s)
				} else {
					headers.Set(HeaderRequestCookie, s)
				}
			}
		}
	}
}

var cookieNameSanitizer = strings.NewReplacer("\n", "-", "\r", "-")

func sanitizeCookieName(n string) string {
	return cookieNameSanitizer.Replace(n)
}

func sanitizeCookieValue(v string) string {
	v = sanitizeOrWarn("Cookie.Value", validCookieValueByte, v)
	if len(v) == 0 {
		return v
	}
	if strings.ContainsAny(v, " ,") {
		return `"` + v + `"`
	}
	return v
}

func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}

func sanitizeOrWarn(fieldName string, valid func(byte) bool, v string) string {
	ok := true
	for i := 0; i < len(v); i++ {
		if valid(v[i]) {
			continue
		}
		log.Printf("net/http: invalid byte %q in %s; dropping invalid bytes", v[i], fieldName)
		ok = false
		break
	}
	if ok {
		return v
	}
	buf := make([]byte, 0, len(v))
	for i := 0; i < len(v); i++ {
		if b := v[i]; valid(b) {
			buf = append(buf, b)
		}
	}
	return string(buf)
}
