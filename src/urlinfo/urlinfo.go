package urlinfo

// URLDB is a simple database of URLs that are malware
type URLDB struct {
	DB map[string]bool
}

// Lookup looks a URL up in a database and returns true if there is a hit (URL is Malware and should be avoided).
func (udb *URLDB) Lookup(url string) bool {
	return udb.DB[url]
}
