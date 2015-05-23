package main

import "github.com/sanbornm/go-selfupdate/selfupdate"

func init() {
	var updater = &selfupdate.Updater{
		CurrentVersion: version,
		ApiURL:         "http://update.luzifer.io/",
		BinURL:         "http://update.luzifer.io/",
		DiffURL:        "http://update.luzifer.io/",
		Dir:            "/home/pi/.radiopi.update/",
		CmdName:        "radiopi", // app name
	}

	if updater != nil {
		go updater.BackgroundRun()
	}
}
