package icecast // import "github.com/Luzifer/radiopi/icecast"

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

// Directory represents the XIPH IceCast directory
type Directory struct {
	cacheFile string
	XMLName   xml.Name `xml:"directory"`
	Entries   []Entry  `xml:"entry"`
}

// Entry is one entry of the IceCast directory
// <entry>
//   <server_name>Vox Noctem</server_name>
//   <listen_url>http://r2d2.voxnoctem.com:8000/voxnoctem.mp3</listen_url>
//   <server_type>audio/mpeg</server_type>
//   <bitrate>192</bitrate>
//   <channels>0</channels>
//   <samplerate>0</samplerate>
//   <genre>
//     80s medieval punk industrial gothic goth electro ebm darkwave
//   </genre>
//   <current_song>Infernosounds - Creature Of The Night</current_song>
// </entry>
type Entry struct {
	ServerName  string `xml:"server_name" json:"server_name"`
	ListenURL   string `xml:"listen_url" json:"listen_url"`
	ServerType  string `xml:"server_type" json:"server_type"`
	Bitrate     string `xml:"bitrate" json:"bitrate"`
	Channels    int    `xml:"channels" json:"channels"`
	SampleRate  int    `xml:"samplerate" json:"samplerate"`
	Genre       string `xml:"genre" json:"genre"`
	CurrentSong string `xml:"current_song" json:"current_song"`
}

type byServerName []Entry

func (b byServerName) Len() int           { return len(b) }
func (b byServerName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byServerName) Less(i, j int) bool { return b[i].ServerName < b[j].ServerName }

// New loads the directory state from disk or network
func New(cacheFile string, typeFilter string) (*Directory, error) {
	xmlData := []byte{}
	if _, err := os.Stat(cacheFile); err == nil {
		xmlData, _ = ioutil.ReadFile(cacheFile)
	}

	if len(xmlData) == 0 {
		req, _ := http.NewRequest("GET", "http://dir.xiph.org/yp.xml", nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		xmlData, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
	}

	tmp := Directory{}
	err := xml.Unmarshal(xmlData, &tmp)
	if err != nil {
		return nil, err
	}

	entries := []Entry{}
	if typeFilter == "" {
		entries = tmp.Entries
	} else {
		for _, e := range tmp.Entries {
			if e.ServerType != typeFilter {
				continue
			}
			entries = append(entries, e)
		}
	}

	return &Directory{
		cacheFile: cacheFile,
		Entries:   entries,
	}, nil
}

// Search does a string matching of genre tags and title against the search string
func (d *Directory) Search(search string) []Entry {
	result := []Entry{}
	searches := strings.Split(strings.ToLower(search), " ")
	for _, e := range d.Entries {
		matches := 0
		for _, s := range searches {
			if strings.Contains(strings.ToLower(strings.Join([]string{e.ServerName, e.Genre}, "::::")), s) {
				matches++
			}
		}
		if len(searches) == matches {
			result = append(result, e)
		}
	}

	sort.Sort(byServerName(result))

	return result
}

// SaveCache stores the filtered list of entries to the cacheFile
func (d *Directory) SaveCache() error {
	out, err := xml.Marshal(d)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(d.cacheFile, out, 0600)
	return err
}
