package main

import (
	"SLB/cfg"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

var config = cfg.Config{}

// var count int

var count map[int]int

const serverMethod = -1

//Server key is -1

type Config struct {
	Servers []string `json:"servers"`
	Routes  []Route  `json:"routes"`
	Port    string   `json:"port`
}

type Route struct {
	Route     string   `json:"route"`
	Endpoints []string `json:"endpoints"`
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
	baseURL := r.URL.Path[1:]
	writeToLog("Basepath: /" + baseURL)
	if len(config.Servers) > 0 {
		server := chooseServer(config.Servers, serverMethod)
		writeToLog("Server: " + server)
		proxy(server, w, r)
	} else if len(config.Routes) > 0 {
		for m := range config.Routes {
			route := config.Routes[m].Route
			bURL := strings.Split(route, "/")[1]
			if baseURL == bURL {
				server := chooseServer(config.Routes[m].Endpoints, m)
				writeToLog("Route: " + server)
				proxy(server, w, r)
			}
		}
	}
}

func chooseServer(servers []string, method int) string {
	count[method] = (count[method] + 1) % len(servers)
	fmt.Println(count[method])
	return servers[count[method]]
}

func writeToLog(message string) {
	logFile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logger := log.New(logFile, "", log.LstdFlags)
	logger.Println(message)
	logFile.Close()
}

func reloadConfig(configFile string, config *cfg.Config) {
	var s string
	fmt.Scanln(&s)
	if s == "-t" {
		*config = cfg.Parse(configFile)
	}
}

func main() {
	count = make(map[int]int)

	var configFile = "./config.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	config = cfg.Parse(configFile)
	// fmt.Println(configFile)
	// data, err := ioutil.ReadFile(configFile)
	// err = json.Unmarshal(data, &config)
	// if err != nil {
	// 	panic(err)
	// }
	// for server := range config.Routes {
	// 	fmt.Println(config.Routes[server])
	// }

	port := ":" + config.Port
	if port == ":" {
		port = port + "8080"
	}

	http.HandleFunc("/", handle)
	writeToLog("Port: " + port)
	writeToLog("Starting server...")
	log.Fatal(http.ListenAndServe(port, nil))
}
