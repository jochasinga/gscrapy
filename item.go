package gscrapy

import (
	"io"

	"golang.org/x/net/html"
)

type Item interface {
	Add(string, *html.Node)
	Del(string)
	Get(string) *html.Node
	Set(string, *html.Node)
	Write(w io.Writer) error
}

type BaseItem map[string]*html.Node
