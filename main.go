package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/cvhariharan/SLB/cfg"
)

var config = cfg.Config{}
var count map[int]int

//Server key is -1
const serverMethod = -1

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
	baseURL = strings.Split(baseURL, "/")[0]
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
	writeToLog("Chose server: " + servers[count[method]])
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

//Could be improved but gets the job done
func reloadConfig(configFile string, config chan cfg.Config, wg *sync.WaitGroup) {
	var s string
	var oldConfig cfg.Config
	var t cfg.Config
	for {
		t = cfg.Parse(configFile)
		fmt.Println(reflect.DeepEqual(t, oldConfig))
		if !reflect.DeepEqual(t, oldConfig) {
			config <- t
			fmt.Println("Reloaded")
			oldConfig = t
		}
		fmt.Scanln(&s)
		if s == "exit" {
			fmt.Println("Closing configChannel")
			close(config)
			wg.Done()
			return
		}

	}
}

func launch(server *http.Server, wg *sync.WaitGroup) {
	writeToLog("Port: " + server.Addr)
	writeToLog("Starting server...")
	handler := http.HandlerFunc(handle)
	server.Handler = handler
	server.ListenAndServe()
	wg.Done()
}

func main() {
	var configFile = "./config.json"
	var server *http.Server
	var wg sync.WaitGroup

	// Adding the reload and exit goroutines
	wg.Add(2)

	count = make(map[int]int)

	configChannel := make(chan cfg.Config)

	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}
	go reloadConfig(configFile, configChannel, &wg)

	go func() {
		for config = range configChannel {
			fmt.Println(config)

			port := ":" + config.Port
			if port == ":" {
				port = port + "8080"
			}
			fmt.Println(server)
			if server != nil {
				writeToLog("Server closing: " + server.Addr)
				fmt.Println("Server closing...")
				server.Close()
			}
			server = &http.Server{
				Addr:         port,
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}
			wg.Add(1)
			go launch(server, &wg)
		}
		fmt.Println("final")
		wg.Done()
	}()

	wg.Wait()
}
