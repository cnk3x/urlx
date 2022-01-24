package proxy

import (
	"github.com/Qingluan/merkur"
	"github.com/cnk3x/urlx"
)

// WithUri 使用代理，支持 ss:// ssr:// vmess:// http:// https:// sock5://
func WithUri(proxy string) urlx.Option {
	return func(r *urlx.Request) error {
		if proxy != "" {
			if client := merkur.NewProxyHttpClient(proxy, 10); client != nil {
				r.UseClient(client)
			}
		}
		return nil
	}
}

// Subscribe 使用订阅来当做代理池
func Subscribe(subscribeUri string) urlx.Option {
	pool := merkur.NewProxyPool(merkur.ParseOrder(subscribeUri)...)
	return func(r *urlx.Request) error {
		if client := pool.GetDialer2().ToHttpClient(10); client != nil {
			r.UseClient(client)
		}
		return nil
	}
}
