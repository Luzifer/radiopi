package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"
)

func playStream(res http.ResponseWriter, r *http.Request) {
	if len(r.FormValue("stream")) > 0 {
		streamChangeChan <- r.FormValue("stream")
		http.Error(res, "OK", http.StatusOK)
		return
	}
	http.Error(res, "Please provide a stream", http.StatusInternalServerError)
}

func getFilteredDirectoryList(res http.ResponseWriter, r *http.Request) {
	res.Header().Add("Content-Type", "application/json")
	search := r.URL.Query().Get("search")
	if len(search) < 3 {
		res.Write([]byte("[]"))
		return
	}

	entries := directory.Search(search)
	if len(entries) > 100 {
		res.Write([]byte("[]"))
		return
	}

	out, err := json.Marshal(entries)
	if err != nil {
		http.Error(res, "An error ocurred", http.StatusInternalServerError)
	}

	res.Write(out)
}

func getFavorites(res http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(favorites)
	if err != nil {
		http.Error(res, "An error ocurred", http.StatusInternalServerError)
	}

	res.Header().Add("Content-Type", "application/json")
	res.Write(body)
}

func serveStatic(res http.ResponseWriter, r *http.Request) {
	p := strings.TrimLeft(r.URL.Path, "/")
	if len(p) == 0 {
		p = "index.html"
	}

	switch path.Ext(p) {
	case ".html":
		res.Header().Add("Content-Type", "text/html")
	case ".js":
		res.Header().Add("Content-Type", "application/javascript")
	default:
		http.Error(res, "Nope.", http.StatusNotFound)
		return
	}

	file, err := Asset(path.Join("frontend", p))
	if err != nil {
		http.Error(res, "Nope.", http.StatusNotFound)
		return
	}

	res.Write(file)
}
