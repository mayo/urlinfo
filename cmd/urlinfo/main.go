package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/purell"
	"github.com/mayo/urlinfo"
)

const (
	//URLPrefix for the service
	URLPrefix    = "/urlinfo/1/"
	urlPrefixLen = len(URLPrefix)
)

// parseCleanURL parses the request URI, parsing out the link that needs to be checked, and doing a light validation and normalization on it. http:// is assumed for normalization, as the incoming links have no protocol specified. Alternatively, this could be stripped off at the end.
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

// handler handles the incoming requests, stripping off the service and version prefix.
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("{\"match\": %t}", ok)))
	})
}

// badRequestError is a helper for returning HTTP errors. Idea borrowed from Go's http package.
func badRequestError(w http.ResponseWriter) {
	httpError(w, "400: Bad Request", http.StatusBadRequest)
}

// httpError handles generic http errors. Idea borrowed from Go's http package.
func httpError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func main() {
	var datafile string

	//Handle command line arguments
	flag.StringVar(&datafile, "datafile", "", "Data file to load")
	flag.Parse()

	// If datafile wasn't specified, bail out
	if datafile == "" {
		fmt.Println("No data file specified.")
		os.Exit(1)
	}

	// Initialize a new URL database
	urlDB := urlinfo.NewByteMapURLDB()

	// Load data
	fmt.Println("urlinfo: Loading malware URLs...")
	err := urlDB.Load(datafile)
	if err != nil {
		fmt.Println("urlinfo: Could not load URLs:", err)
		os.Exit(1)
	}

	// Serve requests
	fmt.Println("urlinfo: Starting server")
	http.HandleFunc(URLPrefix, handler(urlDB))
	http.ListenAndServe(":8080", nil)
}
