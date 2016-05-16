package gscrapy

import (
	"bytes"
	"io"
	"net/http"
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
	opt := &Options{
		Request: &http.Request{
			Host:  "hello.duh",
			Close: true,
		},
		BotName:         "greasybot",
		Contact:         "apology@example.com",
		UserAgentFormat: "%s(%s)",
		Headers: http.Header{
			"Accept": {
				"text/html,application/xhtml+xml," +
					"application/xml;q=0.9,*/*;q=0.8",
			},
			"Accept-Language": {"en"},
		},
	}

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
		t.Errorf("Expect %q. Got %q", expected, result)
	}

	expected = "en"
	result = req.Header.Get("accept-language")
	if strings.Compare(result, expected) != 0 {
		t.Errorf("Expect %q. Got %q", expected, result)
	}

	expectedBool := true
	resultBool := req.Close
	if resultBool != expectedBool {
		t.Errorf("Expect %t. Got %t", expectedBool, resultBool)
	}

	expected = "hello.duh"
	result = req.Host
	if strings.Compare(result, expected) != 0 {
		t.Errorf("Expect %q. Got %q", expected, result)
	}
}

var (
	notFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}
	helloHandler = func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello")
	}
	redirHandler = func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
)

/*
// REVIEW: TestRespGen causes TestParseMethod to FAIL
// and coverage to be disabled.
func TestRespGen(t *testing.T) {
	ts1 := httptest.NewServer(http.HandlerFunc(helloHandler))
	ts2 := httptest.NewServer(http.HandlerFunc(notFoundHandler))
	ts3 := httptest.NewServer(http.HandlerFunc(redirHandler))
	opt, err := NewOptions()
	if err != nil {
		t.Error(err)
	}
	results := respGen([]string{ts1.URL, ts2.URL}, opt)

	// Test type
	if reflect.TypeOf(results) != reflect.TypeOf((<-chan *http.Response)(nil)) {
		t.Errorf("Expect %v. Got %v",
			reflect.TypeOf((<-chan *http.Response)(nil)),
			reflect.TypeOf(results),
		)
	}
	ts1.Close()
	ts2.Close()
	ts3.Close()
}
*/

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
