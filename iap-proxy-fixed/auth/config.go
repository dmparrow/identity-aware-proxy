
package auth

import (
    "encoding/json"
    "os"
    "strings"
    "log"
)

type Config struct {
    Routes []Route `json:"routes"`
}

type Route struct {
    Rule     string   `json:"rule"`
    Service  string   `json:"service"`
    OAuth2   OAuth2   `json:"OAuth2"`
    Upstream Upstream `json:"Upstream"`
    AllowedRoles  []string `json:"allowedRoles,omitempty"`
    AllowedGroups []string `json:"allowedGroups,omitempty"`
}

type OAuth2 struct {
    Provider struct {
        Name         string   `json:"name"`
        Issuer       string   `json:"issuer"`
        Realm        string   `json:"realm"`
        ClientID     string   `json:"clientId"`
        ClientSecret string   `json:"clientSecret"`
        Scopes       []string `json:"scopes"`
    } `json:"provider"`
    AuthResponseHeaders []string `json:"authResponseHeaders"`
}

type Upstream struct {
    PassHostHeader bool      `json:"passHostHeader"`
    Servers        []Backend `json:"servers"`
}

type Backend struct {
    URL string `json:"url"`
}

var routeMap = map[string]Route{}

func LoadConfig(path string) (*Config, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    var cfg Config
    err = json.NewDecoder(f).Decode(&cfg)
    return &cfg, err
}

func GetRoute(host string) *Route {
    normalized := strings.ToLower(strings.Split(host, ":")[0])
    log.Printf("GetRoute lookup for: %s", normalized)

    route, ok := routes[normalized]
    if !ok {
        log.Printf("Route NOT found for: %s", normalized)
        return nil
    }

    log.Printf("Route found for: %s", normalized)
    return &route
}
