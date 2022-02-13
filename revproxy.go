package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var Server *http.Server
var Proxies map[string]*httputil.ReverseProxy = make(map[string]*httputil.ReverseProxy)

func AddProxy(route string, port string) error {
	p, err := newProxy("http://localhost:" + port)
	if err != nil {
		return err
	}
	Proxies[route] = p
	return nil
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
		route = "/" // Must be at one root proxy
	}

	Proxies[route].ServeHTTP(w, r)
}
