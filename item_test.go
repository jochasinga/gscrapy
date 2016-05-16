package gscrapy

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/yhat/scrape"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestItemType(t *testing.T) {
	item, err := NewItem()
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(item) != reflect.TypeOf(Item{}) {
		t.Errorf("Expect %v. Got %v\n",
			reflect.TypeOf((Item)(nil)), reflect.TypeOf(item))
	}
}

var itemKeyTable = []struct {
	n        string
	expected string
}{
	{"title", atom.Title.String()},
	{"meta", atom.Meta.String()},
	{"img", atom.Img.String()},
	{"h1", atom.H1.String()},
	{"p", atom.P.String()},
}

func TestCreateItem(t *testing.T) {
	for _, tt := range itemKeyTable {
		item, err := NewItem(tt.n)
		if err != nil {
			t.Error(err)
		}
		for k := range item {
			if k != tt.expected {
				t.Errorf("Expect %v. Got %v\n", tt.expected, k)
			}
		}
	}
}

func TestCreateItemWithBadKeys(t *testing.T) {
	_, err := NewItem("title", "img", "shirley")
	if err == nil {
		t.Errorf("Expect error. Got nil")
	}
}

var keyNodeTable = []struct {
	key  string
	node *html.Node
}{
	{"h1", &html.Node{Data: "shaun"}},
	{"img", &html.Node{Data: "hello"}},
	{"p", &html.Node{Namespace: "meh"}},
	{"title", &html.Node{Type: html.TextNode}},
	{"meta", &html.Node{Type: html.ElementNode}},
}

func TestItemMethods(t *testing.T) {
	item, err := NewItem()
	if err != nil {
		t.Error(err)
	}
	for _, tt := range keyNodeTable {
		// Add
		item.Add(tt.key, tt.node)
		if item[tt.key][0] != tt.node {
			t.Errorf("Expect item[%q][0] = %v. Got %v\n",
				tt.key, tt.node, item[tt.key][0],
			)
		}
		// Get
		if item.Get(tt.key) != tt.node {
			t.Errorf("Expect %v. Got %v", tt.node, item.Get(tt.key))
		}
		// Set
		item.Set(tt.key, nil)
		if item[tt.key][0] != nil {
			t.Errorf("Expect item[%q][0] = nil. Got %v\n",
				tt.key, item[tt.key][0],
			)
		}
		// Reset
		item.Reset(tt.key)
		if len(item[tt.key]) > 0 {
			t.Errorf("Expect 0. Got %d\n", len(item[tt.key]))
		}

		// Del
		item.Del(tt.key)
		if item[tt.key] != nil {
			t.Errorf("Item with key %q should've been deleted", tt.key)
		}
	}
}

var sampleNodes = []struct {
	html string
	node *html.Node
}{
	{"title", &html.Node{}},
	{"meta", &html.Node{}},
	{"h1", &html.Node{}},
}

func TestItemWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	item, err := NewItem()
	if err != nil {
		t.Error(err)
	}
	for _, tt := range sampleNodes {
		item.Add(tt.html, tt.node)
	}
	err = item.Write(buf)
	if err != nil {
		t.Error(err)
	}
	m := map[string][]string{}
	for _, tt := range sampleNodes {
		m[tt.html] = append(m[tt.html], scrape.Text(tt.node))
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(buf.String(), string(b)+"\n") != 0 {
		t.Errorf("Expect %q. Got %q", buf.String(), string(b)+"\n")
	}
}
