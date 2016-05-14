package gscrapy

import (
	"net/http"
	"time"
)

type Options struct {
	BotName         string
	Website         string
	Email           string
	UserAgentFormat string
	Timeout         time.Duration
	Headers         http.Header
}

func NewOptions(opts ...func(o *Options)) *Options {
	o := &Options{
		BotName:         "greasybot",
		Contact:         "apology@example.com",
		UserAgentFormat: "%q(%q)",
		Headers: http.Header{
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language": "en",
		},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
