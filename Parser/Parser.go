package Parser

import (
	"errors"
	"fmt"
	"github.com/zeebo/bencode"
	"os"
)

// only announceUrl and info are required, and need to be checked.
//Rest may be displayed as additional information in code

type Torrent struct {
	announceUrl  string
	info         torrentInfo
	creationDate int64
	title        string
	comment      string
	urlList      []string
	announceList []string
}

// files and length are mutually exclusive, only one can be used!

type torrentInfo struct {
	name        string
	pieceLength int64
	pieces      string
	collections []string
	files       []torrentFileInfo
	length      int64
}

type torrentFileInfo struct {
	crc32  string
	length int64
	md5    string
	mtime  string
	path   []string
	sha1   string
}

func loadtorrentInfo(info map[string]interface{}) (*torrentInfo, error) {
	// convert collections from an interface slice to a string slice
	interfaceSliceCollections := info["collections"].([]interface{})
	collections := make([]string, len(interfaceSliceCollections), len(interfaceSliceCollections))
	for i := range collections {
		collections[i] = interfaceSliceCollections[i].(string)
	}

	// check if files key is present, if not, length key MUST be present
	_, ok := info["files"]
	if ok {

		// convert files from []interface{} to torrentInfoFile struct
		infoMapSlice := info["files"].([]interface{})
		infoFileSlice := make([]torrentFileInfo, len(infoMapSlice), len(infoMapSlice))
		for i, infoInterface := range infoMapSlice {
			infoFile := infoInterface.(map[string]interface{})

			//convert path from interface{} to []string
			pathString := infoFile["path"].([]interface{})
			infoFilePaths := make([]string, len(pathString), len(pathString))
			for i, path := range pathString {
				infoFilePaths[i] = path.(string)
			}

			newFile := torrentFileInfo{
				crc32:  infoFile["crc32"].(string),
				length: infoFile["length"].(int64),
				md5:    infoFile["md5"].(string),
				mtime:  infoFile["mtime"].(string),
				path:   infoFilePaths,
				sha1:   infoFile["sha1"].(string),
			}
			infoFileSlice[i] = newFile
		}

		t := torrentInfo{
			name:        info["name"].(string),
			pieceLength: info["piece length"].(int64),
			pieces:      info["pieces"].(string),
			collections: collections,
			files:       infoFileSlice,
		}
		return &t, nil
	}
	_, length := info["length"]
	if length {
		t := torrentInfo{
			name:        info["name"].(string),
			pieceLength: info["piece length"].(int64),
			pieces:      info["pieces"].(string),
			collections: collections,
			length:      info["length"].(int64),
		}
		return &t, nil
	}
	return nil, errors.New("torrent has neither files key nor length key")
}

func LoadTorrentData(torrentPath string) (*Torrent, error) {
	reader, err := os.Open(torrentPath)
	defer reader.Close()
	if err != nil {
		return nil, err
	}

	var bencodedtorrentData interface{}
	decoder := bencode.NewDecoder(reader)
	err = decoder.Decode(&bencodedtorrentData)
	if err != nil {
		return nil, err
	}
	torrentData := bencodedtorrentData.(map[string]interface{})

	//convert url-list from interface{} to []string
	urlListString := torrentData["url-list"].([]interface{})
	urlListSlice := make([]string, len(urlListString), len(urlListString))
	for i, path := range urlListString {
		urlListSlice[i] = path.(string)
	}

	//convert announce-list from interface{} to []string
	announceListString := torrentData["announce-list"].([]interface{})
	announceListSlice := make([]string, len(announceListString), len(announceListString))
	for i, slice := range announceListString {
		for _, path := range slice.([]interface{}) {
			announceListSlice[i] = path.(string)
		}
	}
	// create torrentInfo object
	info, err := loadtorrentInfo(torrentData["info"].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	t := Torrent{
		announceUrl:  torrentData["announce"].(string),
		info:         *info,
		creationDate: torrentData["creation date"].(int64),
		title:        torrentData["title"].(string),
		comment:      torrentData["comment"].(string),
		urlList:      urlListSlice,
		announceList: announceListSlice,
	}
	return &t, nil
}

func GetTrackerUrl(torrent Torrent) (string, error) {
	if torrent.announceUrl == "" {
		return "", errors.New("No announce url!")
	}
	return torrent.announceUrl, nil
}

func main() {
	t, err := LoadTorrentData("D:\\Coding\\GoProjects\\Torrent\\Parser\\torrent.torrent")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t.announceUrl, t.info)
}
