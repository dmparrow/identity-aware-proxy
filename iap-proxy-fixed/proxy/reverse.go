
package proxy

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func NewReverseProxy(target string, passHost bool) http.Handler {
    u, err := url.Parse(target)
    if err != nil {
        log.Fatalf("invalid upstream: %s", target)
    }

    proxy := httputil.NewSingleHostReverseProxy(u)
    original := proxy.Director
    proxy.Director = func(req *http.Request) {
        original(req)
        if !passHost {
            req.Host = u.Host
        }
    }

    return proxy
}
