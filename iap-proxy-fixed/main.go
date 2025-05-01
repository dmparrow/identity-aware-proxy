
package main

import (
    "log"
    "net/http"
    "strings"
    "os"
    "os/exec"
    "iap/auth"
    "iap/proxy"
)
func listDir(path string) {
    log.Printf("Listing directory: %s", path)
    out, err := exec.Command("ls", "-l", path).CombinedOutput()
    if err != nil {
        log.Printf("Error listing %s: %v", path, err)
    } else {
        log.Print(string(out))
    }
}

func main() {
    cwd, _ := os.Getwd()
    log.Printf("Working dir: %s", cwd)

    listDir("/app")
    listDir("/app/conf")
    config, err := auth.LoadConfig("/app/conf/config.json")
    if err != nil {
        log.Fatalf("Error loading config: %v", err)
    }

    if err := auth.SetupOAuthClients(config); err != nil {
        log.Fatalf("OAuth setup failed: %v", err)
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/login", auth.HandleLogin)
    mux.HandleFunc("/callback", auth.HandleCallback)

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        host := strings.ToLower(r.Host)
        route := auth.GetRoute(host)
        log.Printf("Incoming Host: %s", r.Host)
        log.Printf("Route: %s", route)
        if route == nil {
            http.Error(w, "route not configured", http.StatusNotFound)
            return
        }

        handler := auth.RequireAuth(route, proxy.NewReverseProxy(route.Upstream.Servers[0].URL, route.Upstream.PassHostHeader))
        handler.ServeHTTP(w, r)
    })

    log.Println("🚀 IAP listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
