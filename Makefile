default: bindata

bindata: coffee
	go-bindata frontend

coffee:
	coffee -c frontend/application.coffee
