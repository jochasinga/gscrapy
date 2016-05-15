gscrapy
=======

A [scrapy](http://scrapy.org/) implementation in Go (close to, not quite).

Description
-----------
Explore the goroutines, channels and [pipelines](https://blog.golang.org/pipelines)
for a high-level, highly concurrent web scraper in Go.

Warning!
--------
This is still a work in progress and a lot of things are changing aggressively.

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
        // Create an item map to store scraped data
        item := gs.NewItem("title", "h1")
        // Create a spider
        sp := &gs.BaseSpider{
                Name: "apologybot",
                Contact: "apology@mail.com",
                // Assign to spider
                Item : item,
                StartURLs: []string{
                        "http://techcrunch.com/",
                        "https://www.reddit.com/"
                        "https://en.wikipedia.org",
                        "https://news.ycombinator.com/",
                        "https://www.buzzfeed.com/",
                        "http://digg.com",
                },
        }

        // Loop over the items channel
        for item := range sp.Crawl() {
                data := map[string][]string{}
                for key, nodes := range item {
                        for node := range nodes {
                                data[key] = append(data[key], scrape.Text(node))
                        }
                }
                jsn, _ := json.Marshal(data)
                fmt.Println(jsn)
        }
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
