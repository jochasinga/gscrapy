package gscrapy

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/yhat/scrape"

	"golang.org/x/net/html"
)

type Item interface {
	Add(string, *html.Node)
	Del(string)
	Get(string) *html.Node
	Write(w io.Writer) error
}

type BaseItem map[string]*html.Node

func NewItem(htmlStr ...string) Item {
	item := &BaseItem{}
	for _, st := range htmlStr {
		item.Add(st, nil)
	}
	return item
}

func (item BaseItem) Add(key string, node *html.Node) {
	capKey := strings.Title(key)
	item[capKey] = node
}

func (item BaseItem) Del(key string) {
	capKey := strings.Title(key)
	delete(item, capKey)
}

func (item BaseItem) Get(key string) *html.Node {
	capKey := strings.Title(key)
	return item[capKey]
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
