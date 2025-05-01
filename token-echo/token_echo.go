
package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintln(w, "<html><body><h1>Token Echo Page</h1><pre>")
    for k, v := range r.Header {
        fmt.Fprintf(w, "%s: %s\n", k, v)
    }
    fmt.Fprintln(w, "</pre></body></html>")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Serving token echo page on :8080")
    http.ListenAndServe(":8080", nil)
}
