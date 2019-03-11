package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var config = Config{}
var count int

type Config struct {
	Servers []string `json:"servers"`
	Port    string   `json:"port`
}

func proxy(target string, w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = url.Host

	proxy.ServeHTTP(w, r)
}

func handle(w http.ResponseWriter, r *http.Request) {
	count = (count + 1) % len(config.Servers)
	proxy(config.Servers[count], w, r)
}

func main() {
	data, err := ioutil.ReadFile("./config.json")
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	for server := range config.Servers {
		fmt.Println(config.Servers[server])
	}

	port := ":" + config.Port
	if port == ":" {
		port = port + "8080"
	}

	http.HandleFunc("/", handle)
	fmt.Println("Port" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
