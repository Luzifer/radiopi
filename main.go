package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"

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
	favoritesFile    *string
	favorites        []favorite
	version          = "dev"
	netMonInterfaces *string
)

func init() {
	var err error

	deadChan = make(chan bool)
	streamChangeChan = make(chan string)

	storeFile = flag.String("cache", "/home/pi/.radiopi", "Cache file to store last stream URL")
	directoryCache := flag.String("directory-file", "/home/pi/.radiopi.directory", "File to cache the IceCast directory to")
	favoritesFile = flag.String("favorites", "/home/pi/.radiopi.favorites", "File to store the favorites in")
	listen = flag.String("listen", ":80", "Listen address for the daemon")
	netMonInterfaces = flag.String("interfaces", "eth0,wlan0", "Interfaces to watch for traffic")
	flag.Parse()

	netmon.Interfaces = strings.Split(*netMonInterfaces, ",")

	directory, err = icecast.New(*directoryCache, "audio/mpeg")
	if err != nil {
		panic(err)
	}
	directory.SaveCache()

	loadFavorites()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/v1/play", playStream).Methods("POST")
	r.HandleFunc("/v1/search", getFilteredDirectoryList).Methods("GET")
	r.HandleFunc("/v1/favorites", getFavorites).Methods("GET")
	r.HandleFunc("/v1/version", getVersion).Methods("GET")
	r.PathPrefix("/").HandlerFunc(serveStatic)

	http.Handle("/", r)
	go http.ListenAndServe(*listen, nil)

	if lastStream, err := ioutil.ReadFile(*storeFile); err == nil {
		playingStream = string(lastStream)
		go restartPlayer()
	}

	netCheck := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-deadChan:
			go restartPlayer()
		case stream := <-streamChangeChan:
			playingStream = stream
			if playerCmd != nil && playerCmd.Process != nil {
				exec.Command("/usr/bin/killall", "mpg123").Run()
			} else {
				deadChan <- true
			}
			ioutil.WriteFile(*storeFile, []byte(stream), 0600)
		case <-netCheck.C:
			expectedRateChange := uint64(64 / 8) // 1s in a 64kbps stream
			if netmon.RateRX < expectedRateChange && playerCmd != nil && playerCmd.Process != nil {
				playerCmd.Process.Kill()
			}
		}
	}
}

func restartPlayer() {
	if playingStream != "off" {
		playerCmd = exec.Command("/usr/bin/mpg123", "-b", "1024", "--no-gapless", playingStream)
		playerCmd.Run()
	} else {
		for playingStream == "off" {
			<-time.After(time.Second)
		}
	}
	deadChan <- true
}
