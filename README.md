# urlinfo

A simple service that keeps in-memory list of URLs and responds to queries asking whether a given URL is present or not. Note that the service does not care about the protocol (http/https) and http will always be assumed internally. Data files with URL lists need to start with `http://`.

`urlinfo` is naive and assumes there is enough memory to load the given data file. While loading, URLs are normalized and hashed to 16 bytes, to save on memory use.

The service will first load the file to memory, and then start service requests.

## Installing

Install the service with `go get github.com/mayo/urlinfo/...`. This will download the package and install `urlinfo` binary in your Go bin directory.

## Start

To read URLs from `urls.txt` file, execute: `urlinfo -datafile urls.txt`

## Checking URLs

To check whether a URL is contained, query the service like so: `http://localhost:8080/urlinfo/1/domain.tlc/path`. This queries the service to check if `http://domain.tlc/path` exists in the URL list:
* If the url exists in the list, the service will respond with `{"match": true}`.
* If the url does not exist in the list, the service will respond with `{"match": false}`.

## Running in Docker

The included script `build-docker-image.sh` will cross-compile the Go binary for Linux and build a minimalistic Docker image named `urlinfo`. It exposes port 8080, and expects the data file to be in `/data/dataset.txt` (`/data` is marked as volume).

An instance then can be started with: `docker run -d -p 8080:8080 -v /path/to/datafile.txt:/data/dataset.txt urlinfo`.