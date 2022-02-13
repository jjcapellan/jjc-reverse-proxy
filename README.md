# JJC-REVERSE-PROXY

Minimal and convenience library to share the same port between multiple services (api, frontend, auth, ...) using reverse proxies.

## How to use
In this example there are 3 servers running in the same machine:
* Frontend server on port 8081
* Api server on port 8082
* Reverse proxy server on port 8080 (the only port accessible from internet)  
The reverse proxy redirects all request to local ports 8081 or 8082 depending on its routes.  

```golang
import(
    "sync"
    proxy "github.com/jjcapellan/jjc-reverse-proxy"    
)

func main(){
    wg := &sync.WaitGroup{}
    
    // Api server is running and listening on port 8082.
    // All request to routes with path prefix "api" (Example: www.site.com/api/user/23415)
    // will be redirected to port 8082.
    proxy.AddProxy("api", "8082")
    // Frontend server is running and listening on port 8081.
    // Must be at least one proxy configured with this path "\". All request not redirected
    // by other proxies will be redirected to port 8081.
    proxy.AddProxy("/", "8081")

    go func(){
        defer wg.Done()
        wg.Add(1)
        // Starts the main server on public port 8080. All request are received in this port
        // and redirected to the api port or frontend port according to its routes.
        proxy.Start("8080")
        defer proxy.Server.Close()
    }

    wg.Wait()
}
```

## Dependencies
This library is built over standard golang libraries, so it hasn't external dependencies.
