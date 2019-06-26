package urlinfo

import (
	"bufio"
	"hash/fnv"
	"os"
	"strings"
)

// URLDB is a generic interface for lookup and loading a URL database
type URLDB interface {
	Lookup(url string) bool
	Load(filename string) error
	Add(url string)
}

// StringMapURLDB is a map based URL database, storing the URL (key) as string.
type StringMapURLDB struct {
	db map[string]bool
}

// NewStringMapURLDB creates a new instance of MapURLDB with an empty map
func NewStringMapURLDB() StringMapURLDB {
	mdb := StringMapURLDB{}
	mdb.db = make(map[string]bool)

	return mdb
}

// Lookup given URL in data store and return true if the URL is present
func (mdb StringMapURLDB) Lookup(url string) bool {
	_, ok := mdb.db[url]
	return ok
}

// Add a new entry to the DB
func (mdb StringMapURLDB) Add(url string) {
	mdb.db[url] = true
}

// Load data into the internal map. The file is expected to have a normalized url per line, starting with http://
func (mdb StringMapURLDB) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		mdb.Add(url)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// ByteSum is a 16 byte array
type ByteSum [16]byte

// ByteSumBoolMap is maps ByteSum to boolean
type ByteSumBoolMap map[ByteSum]bool

// ByteMapURLDB is a map based URL database, storing a hashed URL
type ByteMapURLDB struct {
	db ByteSumBoolMap
}

// NewByteMapURLDB initiqlized a new ByteMapURLDB with an empty map
func NewByteMapURLDB() ByteMapURLDB {
	hmdb := ByteMapURLDB{}
	hmdb.db = make(ByteSumBoolMap)

	return hmdb
}

// Hash the given string (URL)
func (hmdb ByteMapURLDB) Hash(data string) (out ByteSum) {
	h := fnv.New64a()
	h.Write([]byte(data))
	copy(out[:], h.Sum(nil))
	return
}

// Lookup given URL in data store and return true if the URL is present
func (hmdb ByteMapURLDB) Lookup(url string) bool {
	_, ok := hmdb.db[hmdb.Hash(url)]
	return ok
}

// Add a new entry to the DB
func (hmdb ByteMapURLDB) Add(url string) {
	hmdb.db[hmdb.Hash(url)] = true
}

// Load data into the internal map. The file is expected to have a normalized url per line, starting with http://.
func (hmdb ByteMapURLDB) Load(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		hmdb.Add(url)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
