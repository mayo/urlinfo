package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mayo/urlinfo"
)

var testValidURLSet = []struct {
	input    string
	expected string
}{
	{"foo", "http://foo"},
	{"domain.com", "http://domain.com"},
	{"host:80", "http://host"},
	{"host:/file", "http://host/file"},
	{"host:123/file", "http://host:123/file"},
	{"domain.com/path/", "http://domain.com/path/"},
	{"domain.com/path/file.html", "http://domain.com/path/file.html"},
	{"domain.com/path/file?query", "http://domain.com/path/file?query="},
	{"domain.com/path?query", "http://domain.com/path?query="},
	{"domain.com/path/file?query=x/y", "http://domain.com/path/file?query=x/y"},
	{"domain.com/path/file?query=x/y#x", "http://domain.com/path/file?query=x/y"},
	{"host:80/#asd/", "http://host/"},
}

var testInvalidURLSet = map[string]int{
	"":            http.StatusBadRequest,
	"/":           http.StatusBadRequest,
	"/index.html": http.StatusBadRequest,
}

var malwareURLs = []string{
	"http://evilfoo.com",
	"http://malware.com",
	"http://foo.com/evil",
}

type Resp struct {
	Malware bool `json:"malware"`
}

func TestParseURLValid(t *testing.T) {
	for _, test := range testValidURLSet {
		// u, err := url.Parse(URLPrefix + test.input)
		// if err != nil {
		// 	t.Fatalf("Could not parse test input URL")
		// }
		u := URLPrefix + test.input

		parsed, err := parseCleanURL(u)
		if err != nil {
			t.Errorf("Could not parse URL: %s\n%s", test.input, err.Error())
		}

		if parsed != test.expected {
			t.Errorf("Expected URL: %s\n  Parsed URL: %s", test.expected, parsed)
		}
	}

}

func TestParseURLInvalid(t *testing.T) {
	for reqURL := range testInvalidURLSet {
		// u, err := url.Parse(URLPrefix + reqURL)
		// if err != nil {
		// 	t.Fatalf("Could not parse test input URL")
		// }
		u := URLPrefix + reqURL

		parsed, err := parseCleanURL(u)
		if err == nil {
			t.Errorf("Expected error when parsing URL: \"%s\"", reqURL)

			if parsed != "" {
				t.Errorf("Parsed URL was not empty on error for URL: \"%s\", was: \"%s\"", reqURL, parsed)
			}
		}
	}
}

func testHandlerQuery(name string, handler http.HandlerFunc, url string, expected bool, t *testing.T) {
	t.Run(name, func(t *testing.T) {
		reqURL := URLPrefix + url

		req, _ := http.NewRequest("GET", reqURL, nil)
		req.RequestURI = reqURL
		res := httptest.NewRecorder()

		handler(res, req)

		if res.Code != http.StatusOK {
			t.Error()
		}

		jResp := Resp{}
		err := json.Unmarshal(res.Body.Bytes(), &jResp)
		if err != nil {
			t.Errorf("Could not unmarshal response")
		}

		if jResp.Malware != expected {
			t.Error()
		}
	})
}

func TestHandler(t *testing.T) {
	urlDB := urlinfo.NewStringMapURLDB()
	for _, url := range malwareURLs {
		urlDB.Add(url)
	}

	handlerFunc := handler(&urlDB)

	testHandlerQuery("existing", handlerFunc, "evilfoo.com", true, t)
	testHandlerQuery("non-existing", handlerFunc, "miss", false, t)
}
