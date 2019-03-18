# SLB - Simple Load-Balancer

SLB was made for simplicity. It was not created to replace Nginx. It came out of a basic need to have a simple load-balancer, that is easy to configure and use for personal projects. Written entirely in Go, SLB currently has only round-robin routing but a ping-based router is also in the works. 

## Configuration
**Basic Config**

```
{
    "servers": 
    [
    "http://192.164.5.2:8080",
    "https://example.com",
    "http://192.164.5.2:9000"
    ],
    "port" : "9000"
}
```

**Standard Config**
```
{
    "routes":
        [
            {
                "route" : "/test",
                "endpoints" : [
                    "http://192.165.33.22:8080",
                    "http://192.165.33.22:8081",
                    "http://192.165.33.22:8082"
                ]
            },
            {
                "route" : "/run",
                "endpoints" : [
                    "http://192.133.42.3:9000",
                    "http://192.133.42.3:9001",
                    "http://192.133.42.3:9002",
                    "http://192.133.42.3:9003"
                ]
            }
        ],
    "port" : "8000"
}
```
*Routes* allows you to specify a set of servers for handling each base address.  
**If the config file has a *servers* key, the *routes* key will be ignored.**  
If the config file is changed, press enter to reload the server. *This is not a graceful restart.*

### Run
```
go run main.go ./conf.json
```

## Features
- Configurable
- Simple to use
- Highly reliable
- Multi route support

### TODO
- [ ] Ping based routing
- [ ] Graceful shutdown
- [ ] Support for caching

Feel free to open issues and send pull-requests.