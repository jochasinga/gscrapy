package gscrapy

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
)

type Spider interface {
	Parse(*http.Response) *Item
}

type BaseSpider struct {
	Name      string
	StartURLs []string
	Options   *Options
	callback  func(*http.Response) *Item
}

func formatUserAgentStr(format string, args ...string) string {
	return fmt.Sprintf(format, args...)
}

func prepRequest(method, url string, opt *Options, ropts []func(*http.Request)) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for _, ropt := range ropts {
		ropt(req)
	}
	if len(opt.BotName) > 0 && len(opt.Contact) {
		req.Header.Set("user-agent", formatStr(
			opt.UserAgentFormat, opt.BotName, opt.Contact))

	}
	return req, nil
}

func respGen(urls []string, opt *Options, ropts []func(*http.Request)) <-chan *http.Response {
	_ = runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	out := make(chan *http.Response)
	wg.Add(len(urls))
	for _, url := range urls {
		go func(url string) {
			req, err := prepRequest("GET", url, opt, ropts)
			if err != nil {
				log.Fataln(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fataln(err)
			}
			out <- resp
			wg.done()
		}(url)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (sp *BaseSpider) Crawl(urls []string, opt *Options, ropts ...func(r *http.Request)) {
	options := NewOptions()
	if len(opt.BotName) > 0 {
		options.BotName = opt.BotName
	}
	if opt.Timeout != 0 {
		options.Timeout = opt.Timeout
	}
	if len(opt.Headers) > 0 {
		options.Headers = opt.Headers
	}
	rc := respGen(urls, opt, ropts)

}

func (sp *BaseSpider) Parse(r *http.Response) *Item {
	return nil
}
