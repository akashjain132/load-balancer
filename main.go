package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var severCount = 0

// These constant is used to define server
const (
	SERVER1 = "https://www.google.com"
	SERVER2 = "https://www.facebook.com"
	SERVER3 = "https://www.yahoo.com"
	PORT    = "1338"
)

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Log the typeform payload and redirect url
func logRequestPayload(proxyURL string) {
	log.Printf("proxy_url: %s\n", proxyURL)
}

// Balance returns one of the servers based using round-robin algorithm
func getProxyURL() string {
	var servers = []string{SERVER1, SERVER2, SERVER3}

	server := servers[severCount]
	severCount++

	// reset the counter and start from the beginning
	if severCount >= len(servers) {
		severCount = 0
	}

	return server
}

// Given a request send it to the appropriate url
func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
	url := getProxyURL()

	logRequestPayload(url)

	serveReverseProxy(url, res, req)
}

func main() {
	// start server
	http.HandleFunc("/", handleRequestAndRedirect)

	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
