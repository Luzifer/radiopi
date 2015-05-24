package main

import (
	"io/ioutil"
	"sort"

	"gopkg.in/yaml.v2"
)

type favorite struct {
	Name string `yaml:"name" json:"name"`
	URL  string `yaml:"url" json:"url"`
}

type byfavoriteName []favorite

func (b byfavoriteName) Len() int           { return len(b) }
func (b byfavoriteName) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byfavoriteName) Less(i, j int) bool { return b[i].Name < b[j].Name }

func loadFavorites() {
	res := []favorite{}
	body, err := ioutil.ReadFile(*favoritesFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(body, &res)
	if err != nil {
		return
	}

	sort.Sort(byfavoriteName(res))
	favorites = res
}

func saveFavorites() {
	body, err := yaml.Marshal(favorites)
	if err != nil {
		return
	}

	_ = ioutil.WriteFile(*favoritesFile, body, 0644)
}
