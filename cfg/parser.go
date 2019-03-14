package cfg

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Servers []string `json:"servers"`
	Routes  []Route  `json:"routes"`
	Port    string   `json:"port`
}

type Route struct {
	Route     string   `json:"route"`
	Endpoints []string `json:"endpoints"`
}

func Parse(configFile string) Config {
	var config = Config{}
	data, err := ioutil.ReadFile(configFile)
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}
