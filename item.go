package gscrapy

import (
	"encoding/json"
	"io"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Item interface {
	Add(string, *html.Node)
	Del(string)
	Get(string) *html.Node
	Write(w io.Writer) error
}

type BaseItem map[string]*html.Node

//type BaseItem map[string]scrape.Matcher
//type Item map[string]*Field

func NewItem(htmlStr ...string) Item {
	item := BaseItem{}
	for _, st := range htmlStr {
		item.Add(st, nil)
	}
	return item
}

func (item BaseItem) Add(key string, node *html.Node) {
	a := atom.Lookup([]byte(key))
	if a != 0 {
		item[key] = node
	}
}

func (item BaseItem) Del(key string) {
	delete(item, key)
}

func (item BaseItem) Get(key string) *html.Node {
	return item[key]
}

func (item BaseItem) Write(w io.Writer) error {
	newMap := map[string]string{}
	for key, node := range item {
		newMap[key] = scrape.Text(node)
	}
	err := json.NewEncoder(w).Encode(newMap)
	if err != nil {
		return err
	}
	return nil
}
