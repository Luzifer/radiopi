package main

import "github.com/Luzifer/gobuilder/autoupdate"

func init() {
	updater := autoupdate.New("github.com/Luzifer/radiopi", "master")
	updater.SelfRestart = true
	if version != "dev" {
		go updater.Run()
	}
}
