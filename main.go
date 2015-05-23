package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/Luzifer/radiopi/icecast"
	"github.com/gorilla/mux"
)

var (
	playerCmd        *exec.Cmd
	deadChan         chan bool
	streamChangeChan chan string
	playingStream    string
	storeFile        *string
	directory        *icecast.Directory
	listen           *string
	version          = "0.4.0"
)

func init() {
	var err error

	deadChan = make(chan bool)
	streamChangeChan = make(chan string)

	storeFile = flag.String("cache", "/home/pi/.radiopi", "Cache file to store last stream URL")
	directoryCache := flag.String("directory-file", "/home/pi/.radiopi.directory", "File to cache the IceCast directory to")
	listen = flag.String("listen", ":80", "Listen address for the daemon")
	flag.Parse()

	directory, err = icecast.New(*directoryCache, "audio/mpeg")
	if err != nil {
		panic(err)
	}
	directory.SaveCache()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/play", playStream).Methods("POST")
	r.HandleFunc("/v1/search", getFilteredDirectoryList).Methods("GET")
	r.PathPrefix("/").HandlerFunc(serveStatic)

	http.Handle("/", r)
	go http.ListenAndServe(*listen, nil)

	if lastStream, err := ioutil.ReadFile(*storeFile); err == nil {
		playingStream = string(lastStream)
		go restartPlayer()
	}

	for {
		select {
		case <-deadChan:
			go restartPlayer()
		case stream := <-streamChangeChan:
			playingStream = stream
			if playerCmd != nil && playerCmd.Process != nil {
				playerCmd.Process.Kill()
			} else {
				deadChan <- true
			}
			ioutil.WriteFile(*storeFile, []byte(stream), 0600)
		}
	}
}

func restartPlayer() {
	playerCmd = exec.Command("/usr/bin/mpg123", "-b", "1024", "--no-gapless", playingStream)
	playerCmd.Run()
	deadChan <- true
}
