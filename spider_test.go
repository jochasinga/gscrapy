package gscrapy

import (
	"bytes"
	"reflect"
	"strings"
	"sync"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func TestEmptySpiderConstructor(t *testing.T) {
	sp, err := NewSpider()
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(sp) != reflect.TypeOf((*BaseSpider)(nil)) {
		t.Errorf("Expect %v. Got %v",
			reflect.TypeOf((*BaseSpider)(nil)),
			reflect.TypeOf(sp),
		)
	}
	if strings.Compare(sp.Name, "greasybot") != 0 {
		t.Errorf("Expect %q. Got %q", "greasybot", sp.Name)
	}
	opt := NewOptions()
	if !(reflect.DeepEqual(sp.Options, opt)) {
		t.Errorf("Expect %v. Got %v", opt, sp.Options)
	}
	if len(sp.StartURLs) != 0 {
		t.Errorf("Expect %d. Got %d", 0, len(sp.StartURLs))
	}
}

func TestBasicSpiderConstructor(t *testing.T) {
	sp, err := NewSpider(Basic)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(sp) != reflect.TypeOf((*BaseSpider)(nil)) {
		t.Errorf("Expect %v. Got %v",
			reflect.TypeOf((*BaseSpider)(nil)),
			reflect.TypeOf(sp),
		)
	}
	if strings.Compare(sp.Name, "greasybot") != 0 {
		t.Errorf("Expect %q. Got %q", "greasybot", sp.Name)
	}
	opt := NewOptions()
	if !(reflect.DeepEqual(sp.Options, opt)) {
		t.Errorf("Expect %v. Got %v", opt, sp.Options)
	}
	if len(sp.StartURLs) != 0 {
		t.Errorf("Expect %d. Got %d", 0, len(sp.StartURLs))
	}
}

func TestPrepRequest(t *testing.T) {
	opt := NewOptions()
	req, err := prepRequest("GET", "http://example.com", opt)
	if err != nil {
		t.Error(err)
	}
	expected := "greasybot(apology@example.com)"
	result := req.Header.Get("user-agent")
	if strings.Compare(result, expected) != 0 {
		t.Errorf("Expect %q. Got %q", expected, result)
	}
	expected = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	result = req.Header.Get("accept")
	if strings.Compare(result, expected) != 0 {
		t.Errorf("Expect %s. Got %s", expected, result)
	}
	expected = "en"
	result = req.Header.Get("accept-language")
	if strings.Compare(result, expected) != 0 {
		t.Errorf("Expect %s. Got %s", expected, result)
	}
}

var parseTable = []struct {
	a atom.Atom
	d string
	t html.NodeType
}{
	{atom.Title, "Great Title", html.ElementNode},
	{atom.Meta, "Info", html.ElementNode},
	{atom.H1, "Big Heading", html.ElementNode},
}

func TestParseMethod(t *testing.T) {
	var wg sync.WaitGroup
	in := make(chan *html.Node)
	for _, tt := range parseTable {
		wg.Add(1)
		go func(a atom.Atom, d string, t html.NodeType) {
			in <- &html.Node{
				DataAtom: a,
				Data:     d,
				Type:     t,
			}
			wg.Done()
		}(tt.a, tt.d, tt.t)
	}
	go func() {
		wg.Wait()
		close(in)
	}()
	sp, err := NewSpider()
	if err != nil {
		t.Error(err)
	}
	items := sp.Parse(in)
	// Test type
	if reflect.TypeOf(items) != reflect.TypeOf((<-chan Item)(nil)) {
		t.Errorf("Expect %v. Got %v",
			reflect.TypeOf((<-chan Item)(nil)),
			reflect.TypeOf(items),
		)
	}
	// Test items
	for item := range items {
		if len(item) != 3 {
			t.Errorf("Expect 3. Got %d", len(item))
		}
	}
}

func TestWriteMethod(t *testing.T) {
	sp, err := NewSpider()
	if err != nil {
		t.Error(err)
	}
	buf := new(bytes.Buffer)
	// Called before Parse should return an error.
	if err = sp.Write(buf); err == nil {
		t.Error("Expect empty items error. Got nil")
	}
}
