package auth

import (
    "net/http"
    "fmt"
    "strings"
    "log"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "golang.org/x/oauth2"
)

func generatePKCEVerifier() (verifier, challenge string) {
    buf := make([]byte, 32)
    _, err := rand.Read(buf)
    if err != nil {
        panic("unable to generate PKCE verifier")
    }

    verifier = base64.RawURLEncoding.EncodeToString(buf)
    hash := sha256.Sum256([]byte(verifier))
    challenge = base64.RawURLEncoding.EncodeToString(hash[:])
    return
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
    host := strings.ToLower(r.Host)
    client, ok := clients[host]
    if !ok {
        http.Error(w, "unauthorized host", http.StatusUnauthorized)
        return
    }

    state := "example-state" // TODO: random + validated
    verifier, challenge := generatePKCEVerifier()

    // Store verifier in a cookie (or session store if implemented)
    http.SetCookie(w, &http.Cookie{
        Name:     "pkce_verifier",
        Value:    verifier,
        Path:     "/",
        HttpOnly: true,
        Secure:   IsProd(),            // 🔐 Only secure in prod
        SameSite: http.SameSiteLaxMode,
    })

    authURL := client.Config.AuthCodeURL(state,
        oauth2.SetAuthURLParam("code_challenge", challenge),
        oauth2.SetAuthURLParam("code_challenge_method", "S256"),
    )

    http.Redirect(w, r, authURL, http.StatusFound)
}



func HandleCallback(w http.ResponseWriter, r *http.Request) {
    host := strings.ToLower(r.Host)
    client, ok := clients[host]
    if !ok {
        http.Error(w, "unauthorized host", http.StatusUnauthorized)
        return
    }

    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "missing code", http.StatusBadRequest)
        return
    }
    verifierCookie, err := r.Cookie("pkce_verifier")
    if err != nil {
        http.Error(w, "missing pkce verifier", http.StatusBadRequest)
        return
    }

    ctx := r.Context()
    token, err := client.Config.Exchange(ctx, code,
        oauth2.SetAuthURLParam("code_verifier", verifierCookie.Value),
    )
    if err != nil {
        http.Error(w, fmt.Sprintf("token exchange failed: %v", err), http.StatusInternalServerError)
        return
    }

    rawIDToken, ok := token.Extra("id_token").(string)
    if !ok {
        http.Error(w, "missing id_token", http.StatusInternalServerError)
        return
    }

    _, err = client.Verifier.Verify(ctx, rawIDToken)
    if err != nil {
        http.Error(w, fmt.Sprintf("invalid ID token: %v", err), http.StatusUnauthorized)
        return
    }

    // In a real system: set a secure session cookie here
     // ✅ Store session (naive example using cookie)
     http.SetCookie(w, &http.Cookie{
        Name:     SessionCookieName,
        Value:    rawIDToken,
        Path:     "/",
        HttpOnly: true,
        // Secure: true, // enable in prod
    })


    http.Redirect(w, r, "/", http.StatusFound)
}
