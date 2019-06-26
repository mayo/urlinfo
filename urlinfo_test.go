package urlinfo_test

import (
	"math/rand"
	"testing"

	"github.com/mayo/urlinfo"
)

const (
	chars        = "abcdefghijklmnopqrstuvwxyz1234567890:/.?=&%"
	benchDBSize  = 1024
	urlCharLimit = 2000
)

var malwareURLs = map[string]bool{
	"evilfoo.com":  true,
	"malware.com":  true,
	"foo.com/evil": true,
}

func TestStringMapLookupMiss(t *testing.T) {
	urlDB := urlinfo.StringMapURLDB{DB: malwareURLs}

	if ok := urlDB.Lookup("foo.com"); ok {
		t.Fatal()
	}
}

func TestStringMapLookupHit(t *testing.T) {
	urlDB := urlinfo.StringMapURLDB{DB: malwareURLs}

	if ok := urlDB.Lookup("malware.com"); !ok {
		t.Fatal()
	}
}

func benchmarkLookup(db urlinfo.URLDB, keys []string, b *testing.B) {
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, k := range keys {
			db.Lookup(k)
		}
	}
}

func generateKey(keyLen int) string {
	// Always make it at least 1 character in length
	key := make([]byte, keyLen+1)

	for i := range key {
		key[i] = chars[rand.Intn(len(chars))]
	}

	return string(key)
}

func BenchmarkByteMap(b *testing.B) {
	urlDB := urlinfo.NewByteMapURLDB()
	keys := make([]string, 0, benchDBSize)

	for i := 0; i < benchDBSize; i++ {
		url := generateKey(rand.Intn(urlCharLimit))
		urlDB.Add(url)
		keys = append(keys, url)
	}

	benchmarkLookup(urlDB, keys, b)
}

func BenchmarkMap(b *testing.B) {
	urlDB := urlinfo.NewStringMapURLDB()
	keys := make([]string, 0, benchDBSize)

	for i := 0; i < benchDBSize; i++ {
		url := generateKey(rand.Intn(urlCharLimit))
		urlDB.Add(url)
		keys = append(keys, url)
	}

	benchmarkLookup(urlDB, keys, b)
}
