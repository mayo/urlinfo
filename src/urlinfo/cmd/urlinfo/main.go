package main

import (
	"errors"
	"fmt"
	"net/http"
	"urlinfo"

	"github.com/PuerkitoBio/purell"
)

const (
	//URLPrefix for the service
	URLPrefix    = "/urlinfo/1/"
	urlPrefixLen = len(URLPrefix)
)

// parseCleanURL parses the request URI, parsing out the link that needs to be checked, and doing a light validation and normalization on it
func parseCleanURL(requestURI string) (cleanURL string, err error) {
	requestURI = requestURI[urlPrefixLen:]

	if len(requestURI) == 0 || requestURI[0] == '/' {
		err = errors.New("urlinfo: Invalid URL")
		return
	}

	parsedURL := "http://" + requestURI

	// Normalize
	cleanURL, err = purell.NormalizeURLString(parsedURL, purell.FlagsSafe|purell.FlagRemoveDotSegments|purell.FlagRemoveFragment|purell.FlagRemoveDuplicateSlashes|purell.FlagSortQuery|purell.FlagSortQuery|purell.FlagDecodeOctalHost|purell.FlagDecodeHexHost|purell.FlagRemoveUnnecessaryHostDots|purell.FlagRemoveEmptyPortSeparator)

	if err != nil {
		err = errors.New("urlinfo: Could not normalize URL")
		return
	}

	return
}

func handler(udb urlinfo.URLDB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the requested URL and normalize it
		lookupURL, err := parseCleanURL(r.RequestURI)

		if err != nil {
			badRequestError(w)
			fmt.Println("urlinfo: Couldn't parse URL:", lookupURL, err)
			return
		}

		// Lookup and response
		ok := udb.Lookup(lookupURL)

		// Respond
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{\"malware\": %t}", ok)))
	})
}

func badRequestError(w http.ResponseWriter) {
	httpError(w, "400: Bad Request", http.StatusBadRequest)
}

func httpError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	w.Write([]byte(message))
}

var urlDB urlinfo.URLDB

func main() {
	urls := map[string]bool{"malware.com": true}
	urlDB = urlinfo.MapURLDB{DB: urls}

	http.HandleFunc(URLPrefix, handler(urlDB))
	http.ListenAndServe(":8080", nil)
}
