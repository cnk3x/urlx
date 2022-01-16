package html

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Process = func(resp *http.Response, body io.ReadCloser) error

func HtmlQuery(process func(s *goquery.Selection) error) Process {
	return func(resp *http.Response, body io.ReadCloser) error {
		defer body.Close()
		doc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return fmt.Errorf("read as html: %w", err)
		}
		return process(doc.Selection)
	}
}

// Struct 将HTML解析为Struct
func Struct(out any, find string, options ...Options) Process {
	return HtmlQuery(func(s *goquery.Selection) error {
		return BindSelection(s.Find(find), out, options...)
	})
}
