package gscrapy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Spider interface {
	Crawl([]string, *Options, ...func(*http.Request))
	Parse(<-chan *html.Node, Item) <-chan Item
	Write(w io.Writer) error
}

type BaseSpider struct {
	Name           string
	AllowedDomains []string
	StartURLs      []string
	Options        *Options
	Item           Item
}

func prepRequest(method, url string, opt *Options, ropts []func(*http.Request)) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for _, ropt := range ropts {
		ropt(req)
	}
	if opt != nil {
		if len(opt.BotName) > 0 {
			req.Header.Set("user-agent", fmt.Sprintf(
				opt.UserAgentFormat, opt.BotName, opt.Contact))
		}
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
				log.Fatalln(err)
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatalln(err)
			}
			out <- resp
			wg.Done()
		}(url)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func rootGen(in <-chan *http.Response) <-chan *html.Node {
	_ = runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	out := make(chan *html.Node)
	for resp := range in {
		wg.Add(1)
		go func(resp *http.Response) {
			root, err := html.Parse(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			out <- root
			wg.Done()
		}(resp)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (sp *BaseSpider) Parse(in <-chan *html.Node) <-chan Item {
	var wg sync.WaitGroup
	out := make(chan Item)
	for root := range in {
		wg.Add(1)
		go func(r *html.Node) {
			for key := range sp.Item {
				// TODO: Handle case when field = 0
				key := strings.ToLower(key)
				field := atom.Lookup([]byte(key))
				node, ok := scrape.Find(r, scrape.ByTag(field))
				if ok {
					sp.Item.Add(key, node)
				}
				out <- sp.Item
			}
			wg.Done()
		}(root)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (sp *BaseSpider) Crawl(urls []string, opt *Options, ropts ...func(r *http.Request)) <-chan Item {
	items := sp.Parse(rootGen(respGen(urls, opt, ropts)))
	return items
}

/*
func (sp *BaseSpider) Write(w io.Writer) error {
	for item := range sp.items {
		err := item.Write(w)
		if err != nil {
			return err
		}
	}
	return nil
}
*/
