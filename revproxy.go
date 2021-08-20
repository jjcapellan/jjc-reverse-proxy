package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	PortProxy    string
	PortApi      string
	PortFrontend string
	ApiRoute     string
	ApiKey       string
}

var api *httputil.ReverseProxy
var frontend *httputil.ReverseProxy
var config Config
var Server *http.Server

func Setup(c Config) {
	config = c
	createProxies(c.PortFrontend, c.PortApi)
	Server = setupServer(c.PortProxy)
}

func createProxies(portFrontend string, portApi string) {
	var err error
	api, err = newProxy("http://localhost:" + portApi)
	if err != nil {
		log.Fatal("Proxy: Error creating api handler")
	}
	frontend, err = newProxy("http://localhost:" + portFrontend)
	if err != nil {
		log.Fatal("Proxy: Error creating frontend handler")
	}
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

func Start() {
	if Server == nil {
		log.Fatal("Proxy: Server not configured. Use Setup(config) before start server.")
	}
	log.Println("Proxy server listening on port ", config.PortProxy)
	err := Server.ListenAndServe()
	if err != nil {
		e := err.Error()
		log.Println(e)
		if strings.Contains(e, "Server closed") {
			log.Println("Proxy: server closed")
			return
		}
		log.Fatalf("Proxy: Error initializing server: %s", err)
	}
}

func newProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}

func routesHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/") {
		r.URL.Path = "/" + r.URL.Path
	}
	if strings.HasPrefix(r.URL.Path, config.ApiRoute) {
		r.Header.Set("x-api-key", config.ApiKey)
		api.ServeHTTP(w, r)
		return
	}
	frontend.ServeHTTP(w, r)
}
