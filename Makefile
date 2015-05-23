VERSION = $(shell grep 'version' main.go | sed 's/.*"\([^"]*\)"/\1/')

default:

autoupdate:
	GOOS=linux GOARCH=arm go build
	go-selfupdate --platform="linux-arm" -o selfupdate radiopi $(VERSION)
	rm radiopi
	s3cmd sync -P selfupdate/ s3://update.luzifer.io/radiopi/
