package main // tobe Parser

import (
	"fmt"
	"github.com/zeebo/bencode"
	"os"
)

// only announceUrl and info are required, and need to be checked.
//Rest may be displayed as additional information in code

type Torrent struct {
	announceUrl  string
	info         TorrentInfo
	creationDate int
	title        string
	comment      string
	urlList      []string
	announceList []string
}

// files and length are mutually exclusive, only one can be used!

type TorrentInfo struct {
	name        string
	pieceLength int64
	pieces      string
	collections []string
	files       map[string]string
	length      int
}

/*
func new(announceUrl string, info TorrentInfo, creationDate int, title string, comment string,
	urlList []string, announceList []string) (Torrent, error) {

	t := Torrent{announceUrl: announceUrl, info: info, creationDate: creationDate, title: title, comment: comment,
		urlList: urlList, announceList: announceList}

	if announceUrl == "" && info == nil {
		return t, errors.New("missing announceURL and info")
	}

	if announceUrl == "" {
		return t, errors.New("missing announceURL")
	}

	if info == nil {
		return t, errors.New("missing info")
	}

	return t, nil
}
*/

func loadTorrentInfo(info map[string]interface{}) TorrentInfo {
	// convert collections from an interface slice to a string slice
	interfaceSliceCollections := info["collections"].([]interface{})
	collections := make([]string, len(interfaceSliceCollections), len(interfaceSliceCollections))
	for i := range collections {
		collections[i] = interfaceSliceCollections[i].(string)
	}

	// check if files key is present, if not, length key MUST be present
	_, ok := info["files"]
	if ok {
		t := TorrentInfo{
			name:        info["name"].(string),
			pieceLength: info["piece length"].(int64),
			pieces:      info["pieces"].(string),
			collections: collections,
			files:       info["files"].(map[string]string),
		}
		return t
	}
	t := TorrentInfo{
		name:        info["name"].(string),
		pieceLength: info["piece length"].(int64),
		pieces:      info["pieces"].(string),
		collections: collections,
		length:      info["length"].(int),
	}
	return t
}

func loadTorrentData(torrentPath string) Torrent {
	reader, err := os.Open("D:\\Coding\\GoProjects\\Torrent\\Parser\\torrent.torrent")
	defer reader.Close()
	if err != nil {
		fmt.Printf("err is: %s\n", err)
	}

	var bencodedtorrentData interface{}
	decoder := bencode.NewDecoder(reader)
	err = decoder.Decode(&bencodedtorrentData)
	if err != nil {
		fmt.Printf("err2 is: %s\n", err)
	}
	torrentData := bencodedtorrentData.(map[string]interface{})

	// convert all values from default interface to their supposed values (string, map, etc.)
	announceUrl := torrentData["announce"].(string)
	info := loadTorrentInfo(torrentData["info"].(map[string]interface{}))
	creationDate := torrentData["creation date"].(int)
	title := torrentData["title"].(string)
	comment := torrentData["comment"].(string)
	urlList := torrentData["url-list"].([]string)
	announceList := torrentData["announce-list"].([]string)

	fmt.Println(announceUrl, info, creationDate, title, comment, urlList, announceList)
	t := Torrent{
		announceUrl:  torrentData["announce"].(string),
		info:         loadTorrentInfo(torrentData["info"].(map[string]interface{})),
		creationDate: torrentData["creation date"].(int),
		title:        torrentData["title"].(string),
		comment:      torrentData["comment"].(string),
		urlList:      torrentData["url-list"].([]string),
		announceList: torrentData["announce-list"].([]string),
	}
	return t
}

func main() {
	loadTorrentData("x")
}
