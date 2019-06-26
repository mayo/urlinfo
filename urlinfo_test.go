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

	malwareFile    = "testdata/malware_mini.txt"
	malwareFileHit = "http://evilfoo.com"
)

var malwareURLs = []string{
	"evilfoo.com",
	"malware.com",
	"foo.com/evil",
}

func loadURLs(urlDB urlinfo.URLDB, urls []string) {
	for _, url := range urls {
		urlDB.Add(url)
	}
}

func testLookup(urlDB urlinfo.URLDB, hitURL string, t *testing.T) {
	t.Run("miss", func(t *testing.T) {
		if ok := urlDB.Lookup("miss"); ok {
			t.Error()
		}
	})

	t.Run("hit", func(t *testing.T) {
		if ok := urlDB.Lookup(hitURL); !ok {
			t.Error()
		}
	})
}

// String Map tests
func TestStringMapLookup(t *testing.T) {
	urlDB := urlinfo.NewStringMapURLDB()
	loadURLs(urlDB, malwareURLs)
	testLookup(urlDB, malwareURLs[0], t)
}

func TestStringMapLoadValid(t *testing.T) {
	urlDB := urlinfo.NewStringMapURLDB()
	err := urlDB.Load(malwareFile)

	if err != nil {
		t.Fatal()
	}

	testLookup(urlDB, malwareFileHit, t)
}

func TestStringMapLoadInvalid(t *testing.T) {
	urlDB := urlinfo.NewStringMapURLDB()
	err := urlDB.Load("foo")

	if err == nil {
		t.Error()
	}
}

// ByteMap tests

func TestByteMapLookup(t *testing.T) {
	urlDB := urlinfo.NewByteMapURLDB()
	loadURLs(urlDB, malwareURLs)
	testLookup(urlDB, malwareURLs[0], t)
}

func TestByteMapLoadValid(t *testing.T) {
	urlDB := urlinfo.NewByteMapURLDB()
	err := urlDB.Load(malwareFile)

	if err != nil {
		t.Fatal()
	}

	testLookup(urlDB, malwareFileHit, t)
}

func TestByteMapLoadInvalid(t *testing.T) {
	urlDB := urlinfo.NewByteMapURLDB()
	err := urlDB.Load("foo")

	if err == nil {
		t.Error()
	}
}

// Benchmarks

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
