package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// Proxy manager server
var Server *http.Server

// Contains all configured reverse proxies. map[route]proxy
var Proxies map[string]*httputil.ReverseProxy = make(map[string]*httputil.ReverseProxy)

// AddProxy adds and configure one reverse proxy. All request with path prefix "route" will be
// redirected to the given port.
//
// Example: AddProxy("api", "8084") -> redirects www.site.com:8080/api/user/john to port 8084
func AddProxy(route string, port string) error {
	p, err := newProxy("http://localhost:" + port)
	if err != nil {
		return err
	}
	Proxies[route] = p
	return nil
}

// Start inits the proxy manager server on the given port.
func Start(port string) error {
	Server = setupServer(port)
	log.Println("Proxy server listening on port ", port)
	err := Server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func newProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

func routesHandler(w http.ResponseWriter, r *http.Request) {
	if string(r.URL.Path[0]) != "/" {
		r.URL.Path = "/" + r.URL.Path
	}

	route := strings.Split(r.URL.Path, "/")[1]

	if _, exist := Proxies[route]; !exist {
		route = "/" // Must be one root proxy
	}

	Proxies[route].ServeHTTP(w, r)
}

func setupServer(portProxy string) *http.Server {
	http.HandleFunc("/", routesHandler)

	srv := &http.Server{
		Addr:         ":" + portProxy,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return srv
}
