package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"urlinfo"
)

// parseURL parses incoming URL for an URL that needs to be checked.
func parseURL(url string) (string, error) {
	parts := strings.SplitN(url, "/", 4)
	if len(parts) < 4 {
		return "", errors.New("invalid url request")
	}

	return parts[3], nil
}

func handler(udb *urlinfo.URLDB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Build the URL from requested URL. Using r.URL avoids having to strip fragments. Alternative would be to use r.RequestURI and cut off at first hash (#)
		// Not using RawPath, it seems to be empty
		checkURL := r.URL.Path

		if r.URL.RawQuery != "" {
			checkURL += "?" + r.URL.RawQuery
		}

		// Parse the requested URL
		urlNoScheme, err := parseURL(checkURL)
		if err != nil {
			badRequestError(w)
			fmt.Printf("urlinfo: Couldn't parse requested URL: %s\n", r.URL)
		}

		if urlNoScheme == "" {
			badRequestError(w)
			fmt.Println("urlinfo: No URL specified")
		}

		// Lookup and response
		ok := udb.Lookup(urlNoScheme)
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
	urlDB = urlinfo.URLDB{DB: urls}

	http.HandleFunc("/urlinfo/1/", handler(&urlDB))
	http.ListenAndServe(":8080", nil)
}
