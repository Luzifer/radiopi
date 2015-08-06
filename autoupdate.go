package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"syscall"
	"time"

	"github.com/Luzifer/gobuilder/builddb"
)

func init() {
	updater := NewGoBuilderUpdate("github.com/Luzifer/radiopi", "master")
	updater.SelfRestart = true
	if version != "dev" {
		go updater.Run()
	}
}

type GoBuilderUpdate struct {
	UpdateInterval    time.Duration
	SelfRestart       bool
	repository        string
	label             string
	runningFile       string
	currentHash       string
	goBuilderFilename string
}

func NewGoBuilderUpdate(repo, label string) *GoBuilderUpdate {
	filename := fmt.Sprintf("%s_%s_%s-%s",
		path.Base(repo),
		label,
		runtime.GOOS,
		runtime.GOARCH,
	)

	if runtime.GOOS == "windows" {
		filename = filename + ".exe"
	}

	return &GoBuilderUpdate{
		UpdateInterval:    time.Minute * 60,
		SelfRestart:       false,
		repository:        repo,
		label:             label,
		runningFile:       os.Args[0],
		goBuilderFilename: filename,
	}
}

func (g *GoBuilderUpdate) Run() error {
	bin, err := ioutil.ReadFile(g.runningFile)
	if err != nil {
		return err
	}

	g.currentHash = fmt.Sprintf("%x", md5.Sum(bin))

	for {
		liveHash, err := g.getGoBuilderHash()
		if err == nil && liveHash != g.currentHash {
			err := g.updateBinary()
			if err == nil && g.SelfRestart {
				syscall.Exec(os.Args[0], os.Args[1:], []string{})
			}
		}
		<-time.After(g.UpdateInterval)
	}
}

func (g *GoBuilderUpdate) getGoBuilderHash() (string, error) {
	url := fmt.Sprintf("https://gobuilder.me/api/v1/%s/hashes/%s.json",
		g.repository,
		g.label,
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP Status != 200")
	}

	out := builddb.HashDB{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	hashes, ok := out[g.goBuilderFilename]
	if !ok {
		return "", fmt.Errorf("Could not find hashes for %s", g.goBuilderFilename)
	}

	return hashes.MD5, nil
}

func (g *GoBuilderUpdate) updateBinary() error {
	dlURL := fmt.Sprintf("https://gobuilder.me/get/%s/%s",
		g.repository,
		g.goBuilderFilename,
	)

	resp, err := http.Get(dlURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	self, err := os.Create(g.runningFile)
	if err != nil {
		return err
	}
	defer self.Close()

	_, err = io.Copy(self, resp.Body)
	return err
}
