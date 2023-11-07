package main // tobe Tracker

import (
	"Torrent/Parser"
	_ "Torrent/Parser"
	"fmt"
	_ "fmt"
	"log"
	"net/http"
)

// Tracker package will implement funcs to take a torrent type and extract tracker response from it (peers and interval)

// HttpTrackerRequest holds all parameters for a get request to an http torrent tracker
type HttpTrackerRequest struct {
	infoHash   string // 20 byte sha1 string of the bencoded form of the torrent info value
	peerId     string // random length20 string
	ip         string // optional parameter of self ip/dns name
	port       int    // port number we are listening on
	uploaded   int    // The total amount uploaded so far, encoded in base ten ascii.
	downloaded int    // The total amount downloaded so far, encoded in base ten ascii.
	left       int    // number of bytes left to download, encoded in base ten ascii
	event      string // optional key which maps to [started, completed, stopped, empty]
}

// http tracker responds with bencoded dict with three keys

type HttpTrackerResponse struct {
	interval      int // num of seconds to wait between regular re-requests
	peers         []peer
	failureReason string // should only appear if there is a failure in communication with the tracker,
	// in that case the two other keys do not matter
}

type peer struct {
	peerId string
	ip     string // ip or dns name
	port   int
}

// GenerateRequest will generate a request to the tracker using the Torrent struct and return the query ready to send or error
func GenerateRequest(torrent Parser.Torrent) (string, error) {
	announceUrl, err := Parser.GetTrackerUrl(torrent)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	req, err := http.NewRequest("GET", announceUrl, nil)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	fmt.Println(req)
	return nil, nil
}

// for func receiving the tracker response, need to parse peer as either compact or not (bep0003 vs bep0023)
func main() {
	t := Parser.LoadTorrentData("D:\\Coding\\GoProjects\\Torrent\\Parser\\torrent.torrent")
	GenerateRequest(t)
}
