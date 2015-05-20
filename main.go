package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/gorilla/mux"
)

var (
	playerCmd        *exec.Cmd
	deadChan         chan bool
	streamChangeChan chan string
	playingStream    string
	storeFile        *string
)

func init() {
	deadChan = make(chan bool)
	streamChangeChan = make(chan string)

	storeFile = flag.String("cache", "/home/pi/.radiopi", "Cache file to store last stream URL")
	flag.Parse()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/play", playStream).Methods("POST")

	http.Handle("/", r)
	go http.ListenAndServe(":80", nil)

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
	playerCmd = exec.Command("/usr/bin/mpg123", "--no-gapless", playingStream)
	playerCmd.Run()
	deadChan <- true
}

func playStream(res http.ResponseWriter, r *http.Request) {
	if len(r.FormValue("stream")) > 0 {
		streamChangeChan <- r.FormValue("stream")
		http.Error(res, "OK", http.StatusOK)
		return
	}
	http.Error(res, "Please provide a stream", http.StatusInternalServerError)
}
