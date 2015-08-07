package main

import (
	"time"

	"github.com/Luzifer/gobuilder/autoupdate"
)

func init() {
	updater := autoupdate.New("github.com/Luzifer/radiopi", "master")
	updater.SelfRestart = true
	updater.UpdateInterval = 10 * time.Minute
	if version != "dev" {
		go updater.Run()
	}
}
