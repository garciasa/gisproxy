package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

// Proxy struct to store elements we need
type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	logger *zap.Logger
}

// NewProxy init a proxy with a given url
func NewProxy(target string) *Proxy {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	targetURL, err := url.Parse(target)
	if err != nil {
		panic("error parsing url")
	}

	pxy := httputil.NewSingleHostReverseProxy(targetURL)

	proxy := &Proxy{
		target: targetURL,
		proxy:  pxy,
		logger: logger,
	}

	return proxy
}

func handle(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	url := params.Get("q")

	log.Printf("Redirecting %s", url)
	p := NewProxy(url)
	p.logger.Info("Redirecting ",
		zap.String("url", url))
	r.Host = p.target.Host
	p.proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", handle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
