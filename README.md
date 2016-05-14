gscrapy
=======

A [scrapy](http://scrapy.org/) implementation in Go.

Description
-----------
Explore the goroutines, channels and [pipelines](https://blog.golang.org/pipelines)
for a high-level, highly concurrent web scraper in Go.

Disclaimer!
-----------
Still a work in progress. Better wait 'til it's released.

Install
-------
```bash

go get github.com/jochasinga/gscrapy

```

Usage
-----
```go

package main

import (
        gs "github.com/jochasinga/gscrapy"
)

func main() {
        item := gs.NewItem("title", "h1")
        sp := &gs.BaseSpider{
                Item: item,
                Name: "myspider",
                StartURLs: []string{
			"http://techcrunch.com/",
			"https://www.reddit.com/",
			"https://en.wikipedia.org",
			"https://news.ycombinator.com/",
			"https://www.buzzfeed.com/",
			"http://digg.com",
		},
        }
        sp.Crawl()
        _ := sp.Write(os.Stdout)
}

```

results:

```bash

$ go run main.go
{"title":"reddit: the front page of the internet","h1":"Main Page"}
{"title": "BuzzFeed","h1":"Buzzfeed"}
{"title": "TechCrunch - The latest technology news and information on startups","h1":"Gauri Nanda of Toymail"}
{"title": "Digg - What the Internet is talking about right now","h1":"Digg"}
{"title": "Wikipedia, the free encyclopedia", "h1":"Main Page"}
{"title": "Hacker News","h1":""}

```
