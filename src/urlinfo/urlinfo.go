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
}

// MapURLDB is a map based URL database, storing the URL (key) as string.
type MapURLDB struct {
	DB map[string]bool
}

// NewMapURLDB creates a new instance of MapURLDB with an empty map
func NewMapURLDB() MapURLDB {
	mdb := MapURLDB{}
	mdb.DB = make(map[string]bool)

	return mdb
}

// Lookup given URL in data store and return true if the URL is present
func (mdb MapURLDB) Lookup(url string) bool {
	return mdb.DB[url]
}

// Load data into the internal map. The file is expected to have a normalized url per line, starting with http://
func (mdb MapURLDB) Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		mdb.DB[url] = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type ByteSum [16]byte
type ByteSumBoolMap map[ByteSum]bool

// ByteMapURLDB is a map based URL database, storing the URL as FNV-128a hash
type ByteMapURLDB struct {
	DB ByteSumBoolMap
}

// NewButeMapURLDB initiqlized a new ByteMapURLDB with an empty map
func NewByteMapURLDB() ByteMapURLDB {
	hmdb := ByteMapURLDB{}
	hmdb.DB = make(ByteSumBoolMap)

	return hmdb
}

// Hash the given string (URL)
func (hmdb ByteMapURLDB) Hash(data string) (out ByteSum) {
	h := fnv.New128a()
	h.Write([]byte(data))
	copy(out[:], h.Sum(nil))
	return
}

// Lookup given URL in data store and return true if the URL is present
func (hmdb ByteMapURLDB) Lookup(url string) bool {
	return hmdb.DB[hmdb.Hash(url)]
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
		// Store the hashed URL
		hmdb.DB[hmdb.Hash(url)] = true
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
