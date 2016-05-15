package gscrapy

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var item = NewItem()

func TestItemType(t *testing.T) {
	_, ok := item.(Item)
	if ok != true {
		t.Errorf("%v does not implement Item", item)
	}
	if reflect.TypeOf(item) != reflect.TypeOf((BaseItem)(nil)) {
		t.Errorf("%v is not BaseItem", item)
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
		testItem := NewItem(tt.n)
		for k := range testItem.(BaseItem) {
			if !(strings.Contains(k, tt.expected)) {
				t.Error("Item key does not match")
			}
		}
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
	buf := new(bytes.Buffer)
	for _, tt := range keyNodeTable {
		// Add
		item.Add(tt.key, tt.node)
		if item.(BaseItem)[tt.key] != tt.node {
			t.Error("Item key and node does not match")
		}
		// Get
		if item.Get(tt.key) != tt.node {
			t.Error("Item key and node does not match")
		}
		// FIXME: Write
		err := item.Write(buf)
		if err != nil {
			t.Error(err)
		}
		m := map[string]*html.Node{tt.key: tt.node}
		b, err := json.Marshal(m)
		if err != nil {
			t.Error(err)
		}
		if bytes.Compare(buf.Bytes(), b) != 0 {
			t.Log(buf.String(), string(b))
			t.Fail()
		}
		// Del
		item.Del(tt.key)
		if item.(BaseItem)[tt.key] != nil {
			t.Error("Item key should be deleted")
		}

	}
}
