package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	_ "net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

// Config environment

const host string = "http://localhost"
const proxyPort string = "8080"
const apiRoute string = "/api"
const apiKey string = "abcd"

// Fake client
var client *http.Client = &http.Client{}

func checkRoute(req *http.Request, wanted string, t *testing.T) {
	response, _ := client.Do(req)
	body, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	path := req.URL.Path
	strBody := string(body)
	if strBody != wanted {
		t.Errorf("Bad redirection. Request path: %s. Request was sent to: %s", path, strBody)
	}
}

func buildRequest(path string) *http.Request {
	req, _ := http.NewRequest("GET", host+":"+proxyPort+path, nil)
	return req
}

func TestInit(t *testing.T) {

	// Api fake server
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("x-api-key") != apiKey {
			t.Error("Wrong x-api-key header")
		}
		w.Write([]byte("apiServer"))

	}))
	apiUrl, _ := url.Parse(apiServer.URL)
	apiPort := apiUrl.Port()
	defer apiServer.Close()

	// Frontend fake server
	frontendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("frontendServer"))
	}))
	feUrl, _ := url.Parse(frontendServer.URL)
	frontendPort := feUrl.Port()
	defer frontendServer.Close()

	// Reverse proxy

	wg := &sync.WaitGroup{}
	wg.Add(1)
	c := Config{
		PortProxy:    proxyPort,
		PortApi:      apiPort,
		PortFrontend: frontendPort,
		ApiRoute:     apiRoute,
		ApiKey:       apiKey,
	}

	Setup(c)
	proxyServer := Server

	go func() {
		wg.Done()
		Start()
	}()

	// Fake request to api

	apiReq1 := buildRequest("/api/user")
	apiReq2 := buildRequest("/api")
	apiReq3 := buildRequest("/api/")

	// Fake request to frontend

	feReq1 := buildRequest("/page1")
	feReq2 := buildRequest("/")
	feReq3 := buildRequest("/page/api")

	checkRoute(apiReq1, "apiServer", t)
	checkRoute(apiReq2, "apiServer", t)
	checkRoute(apiReq3, "apiServer", t)
	checkRoute(feReq1, "frontendServer", t)
	checkRoute(feReq2, "frontendServer", t)
	checkRoute(feReq3, "frontendServer", t)

	proxyServer.Close()
	wg.Wait()
}
