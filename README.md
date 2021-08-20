# Important: Still in alpha (not use in production)

Basic Reverse Proxy oriented to a local specific environment (proxy, frontend and api servers, all in localhost).  

This is the proxy behavior:
* By default all request are sent to the frontend server.
* Only those routes with a especific prefix (ex: "/api") are sent to the apirest server.
* Adds an "x-api-key" header to all request sent to api server.

## How to use
```golang
import(
    "sync"
    proxy "github.com/jjcapellan/jjc-reverse-proxy"    
)

func main(){
    wg := &sync.WaitGroup{}

    c := proxy.Config{
            PortProxy:    "8080",
            PortApi:      "8081",
            PortFrontend: "8082",
            ApiRoute:     "/api",
            ApiKey:       "myapikey",
        }

    proxy.Setup(c)
    go func(){
        defer wg.Done()
        wg.Add(1)
        proxy.Start()
    }

    wg.Wait()
}
```
