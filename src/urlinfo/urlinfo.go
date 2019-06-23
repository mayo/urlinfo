package urlinfo

type URLDB interface {
	Lookup(url string) bool
}

// MapURLDB is a simple database of URLs that are malware
type MapURLDB struct {
	DB map[string]bool
}

func NewMapURLDB() MapURLDB {
	mdb := MapURLDB{}
	mdb.DB = make(map[string]bool)

	return mdb
}

// Lookup given URL in data store and return true if the URL is present
func (mdb MapURLDB) Lookup(url string) bool {
	return mdb.DB[url]
}
