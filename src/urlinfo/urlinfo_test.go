package urlinfo_test

import (
	"testing"
	"urlinfo"
)

var badURLs = map[string]bool{
	"evilfoo.com":  true,
	"malware.com":  true,
	"foo.com/evil": true,
}

func TestLookupMiss(t *testing.T) {
	urlDB := urlinfo.MapURLDB{DB: badURLs}

	if ok := urlDB.Lookup("foo.com"); ok {
		t.Fatal()
	}
}

func TestLookupHit(t *testing.T) {
	urlDB := urlinfo.MapURLDB{DB: badURLs}

	if ok := urlDB.Lookup("malware.com"); !ok {
		t.Fatal()
	}
}
