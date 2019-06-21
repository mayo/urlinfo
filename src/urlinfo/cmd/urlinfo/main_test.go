package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlinfo"
)

var testValidURLSet = [...]string{
	"foo",
	"domain.com",
	"host:port",
	"host/file",
	"host:port/file",
	"domain.com/path/",
	"domain.com/path/file",
	"domain.com/path/file?query",
	"domain.com/path?query",
	"domain.com/path/file?query=x/y",
}

var testInvalidURLSet = [...]string{
	"/",
	" foo",
	" /file",
}

var badURLs = map[string]bool{
	"evilfoo.com":  true,
	"malware.com":  true,
	"foo.com/evil": true,
}

var baseURL = "/urlinfo/1/"

type Resp struct {
	Malware bool `json:"malware"`
}

func TestParseURLValid(t *testing.T) {
	for _, url := range testValidURLSet {
		reqURL := baseURL + url
		parsed, err := parseURL(reqURL)
		if err != nil {
			t.Fatalf("Could not parse URL: %s\n%s", reqURL, err.Error())
		}

		if parsed != url {
			t.Fatalf("Expected URL: %s\n  Parsed URL: %s", url, parsed)
		}
	}

}

func TestParseURLEmpty(t *testing.T) {
	parsed, err := parseURL(baseURL)
	if err != nil {
		t.Fatal()
	}

	if parsed != "" {
		t.Fatalf("Expected empty string, received: %s", parsed)
	}
}

//TODO: Make this pass
func TestParseURLInvalid(t *testing.T) {
	reqURL := baseURL + "/"
	parsed, err := parseURL(reqURL)
	if err != nil {
		t.Fatal()
	}

	if parsed != "" {
		t.Fatalf("Expected empty string, received: %s", parsed)
	}
}

func TestHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", baseURL+"evilfoo.com", nil)
	res := httptest.NewRecorder()

	urlDB := urlinfo.URLDB{DB: badURLs}
	handler(&urlDB)(res, req)

	if res.Code != http.StatusOK {
		t.Fatal()
	}

	jResp := Resp{}
	err := json.Unmarshal(res.Body.Bytes(), &jResp)
	if err != nil {
		t.Fatal()
	}

	if jResp.Malware != true {
		t.Fatal()
	}
}
