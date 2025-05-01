package auth

import (
    "context"
    "fmt"
    "strings"
    "log"
    "github.com/coreos/go-oidc/v3/oidc"
    "golang.org/x/oauth2"
)

type clientCtx struct {
    Config   *oauth2.Config
    Verifier *oidc.IDTokenVerifier
}

var clients = make(map[string]*clientCtx)
var routes = make(map[string]Route)


func SetupOAuthClients(cfg *Config) error {
    redirectScheme := "http"
    if IsProd() {
        redirectScheme = "https"
    }

    ctx := context.Background()
    for _, route := range cfg.Routes {
        rule := strings.ToLower(route.Rule)
        issuerURL := fmt.Sprintf("%s/realms/%s", route.OAuth2.Provider.Issuer, route.OAuth2.Provider.Realm)
        provider, err := oidc.NewProvider(ctx, issuerURL)
        if err != nil {
            return fmt.Errorf("OIDC setup failed for route %s: %w", rule, err)
        }

        oauthCfg := &oauth2.Config{
            ClientID:     route.OAuth2.Provider.ClientID,
            ClientSecret: route.OAuth2.Provider.ClientSecret,
            RedirectURL:  fmt.Sprintf("%s://%s/callback", redirectScheme, route.Rule),
            Endpoint:     provider.Endpoint(),
            Scopes:       append([]string{oidc.ScopeOpenID}, route.OAuth2.Provider.Scopes...),
        }
        clients[rule] = &clientCtx{
            Config:   oauthCfg,
            Verifier: provider.Verifier(&oidc.Config{ClientID: route.OAuth2.Provider.ClientID}),
        }
        routes[rule] = route
    }
    return nil
}