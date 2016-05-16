package gscrapy

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/yhat/scrape"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Item map[string][]*html.Node

// NewItem create a new item with predefined keys
// from optional htmlStr.
func NewItem(htmlStr ...string) (Item, error) {
	item := Item{}
	if len(htmlStr) > 0 {
		for _, key := range htmlStr {
			key := strings.ToLower(key)
			a := atom.Lookup([]byte(key))
			if a != 0 {
				item.Reset(key)
			} else {
				return nil, errors.New(`
					One of more HTML keys are not compatible.
					See https://godoc.org/golang.org/x/net/html/atom#Atom`)
			}
		}
	}
	return item, nil
}

// Add adds the key, node pair to the item. It appends
// to any existing nodes associated with key.
func (item Item) Add(key string, node *html.Node) {
	lowered := strings.ToLower(key)
	a := atom.Lookup([]byte(lowered))
	if a != 0 {
		item[lowered] = append(item[lowered], node)
	}
}

// Set sets the item entries associated with key to the
// single element node. It replaces any existing nodes
// associated with key.
func (item Item) Set(key string, node *html.Node) {
	lowered := strings.ToLower(key)
	a := atom.Lookup([]byte(lowered))
	if a != 0 {
		item[lowered] = []*html.Node{node}
	}
}

// Reset resets the node slice at the specified key
// to an empty one.
func (item Item) Reset(key string) {
	lowered := strings.ToLower(key)
	a := atom.Lookup([]byte(lowered))
	if a != 0 {
		item[lowered] = []*html.Node{}
	}
}

// Del deletes the entry at the specified key.
// If the key doesn't exist, nothing is deleted.
func (item Item) Del(key string) {
	lowered := strings.ToLower(key)
	delete(item, lowered)
}

// Get gets the first node stored in an item's key.
// To access multiple nodes, access the map directly.
func (item Item) Get(key string) *html.Node {
	lowered := strings.ToLower(key)
	return item[lowered][0]
}

// Write writes the value from the item to writer w as JSON bytes.
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
