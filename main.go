package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(logFile, "", log.LstdFlags)

	var configFile = "./config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	fmt.Println(configFile)
	data, err := ioutil.ReadFile(configFile)
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
	logger.Println("Port: " + port)
	logger.Println("Starting server...")
	logger.Fatal(http.ListenAndServe(port, nil))
}
