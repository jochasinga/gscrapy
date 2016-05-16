gscrapy
=======

A [scrapy](http://scrapy.org/) implementation in Go (close to, not quite).

Description
-----------
I wanted to explore the goroutines, channels and [pipelines](https://blog.golang.org/pipelines)
for a high-level, highly concurrent and easy to use web scraper in Go.

Warning!
--------
This is still a work in progress! Watch for a release.

Usage
-----

```go

package main

import (
        "os"

        gs "github.com/jochasinga/gscrapy"
)

startURLs = []string{
        "http://techcrunch.com/",
        "https://www.reddit.com/"
        "https://en.wikipedia.org",
        "https://news.ycombinator.com/",
        "https://www.buzzfeed.com/",
        "http://digg.com",
},

func main() {
        // Create an item map to store scraped data
        item := gs.NewItem("title", "h1")
        // Create a spider
        sp := &gs.BaseSpider{
                Name: "apologybot",
                Contact: "apology@gmail.com",
                // Assign to spider
                Item : item,
                StartURLs: startURLs,
        }

        // OR create a default spider
        // sp := NewSpider(gs.Basic)
        // sp.Item = item
        // _ = sp.Crawl(startURLs...)

        _ = sp.Crawl()
        sp.Write(os.Stdout)
}

```

results (hopefully):

```bash

$ go run main.go
{"title":["reddit: the front page of the internet"],["h1":"Main Page"]}
{"title":["BuzzFeed","h1":"Buzzfeed"]}
{"title":["TechCrunch - The latest technology news and information on startups"],"h1":["Gauri Nanda of Toymail"]}
{"title": ["Digg - What the Internet is talking about right now"],"h1":["Digg"]}
{"title": ["Wikipedia, the free encyclopedia"], ["h1":"Main Page"]}
{"title": ["Hacker News"],"h1":[""]}

```
