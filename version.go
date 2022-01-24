package urlx

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

const Version = "0.0.1"

var DefaultUserAgent HeaderOption

func init() {
	type x struct{}
	dua := fmt.Sprintf("urlx/%s (%s) golang/%s(%s %s)", Version, reflect.TypeOf(x{}).PkgPath(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
	DefaultUserAgent = UserAgent(dua)
}

// Default 默认的请求器
func Default(ctx context.Context) *Request {
	ms := time.Millisecond
	return New(ctx).HeaderWith(DefaultUserAgent).TryAt(ms*300, ms*800, ms*1500)
}
