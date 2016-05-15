package gscrapy

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/yhat/scrape"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Item map[string][]*html.Node

// NewItem create and returns *Item. htmlStr provided
// are the targeted HTML tags to scrape.
func NewItem(htmlStr ...string) Item {
	item := Item{}
	// Start fresh
	if len(htmlStr) > 0 {
		for _, str := range htmlStr {
			item.Set(str, nil)
		}
	}
	return item
}

func (item Item) Add(key string, node *html.Node) {
	lowered := strings.ToLower(key)
	a := atom.Lookup([]byte(lowered))
	if a != 0 {
		item[lowered] = append(item[lowered], node)
	}
}

func (item Item) Set(key string, node *html.Node) {
	//capped := strings.Title(key)
	lowered := strings.ToLower(key)
	a := atom.Lookup([]byte(lowered))
	if a != 0 {
		item[lowered] = []*html.Node{node}
	}
}

func (item Item) Del(key string) {
	lowered := strings.ToLower(key)
	delete(item, lowered)
}

func (item Item) Get(key string) *html.Node {
	lowered := strings.ToLower(key)
	return item[lowered][0]
}

func (item Item) Write(w io.Writer) error {
	newMap := map[string][]string{}
	for key, nodes := range item {
		for _, node := range nodes {
			text := scrape.Text(node)
			newMap[key] = append(newMap[key], text)
		}
	}
	err := json.NewEncoder(w).Encode(newMap)
	if err != nil {
		return err
	}
	return nil
}
