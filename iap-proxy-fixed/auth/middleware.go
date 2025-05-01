package auth

import (
    "net/http"
    "log"
    "time"
    "strings"
)
var SessionCookieName = "iap_session"

func Middleware(route *Route, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        now := time.Now().UTC().Format(time.RFC3339)

        cookie, err := r.Cookie(SessionCookieName)
        if err != nil || cookie.Value == "" {
            log.Printf("[%s] No session cookie for %s -> Redirecting to /login", now, r.URL.Path)
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

        session, err := GetSession(r)
        
        if err != nil || !session.IsValid() {
            log.Printf("[%s] Unauthenticated access to %s -> Redirecting to login", now, r.URL.Path)
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

        userRoles := strings.Join(session.Roles, ", ")

        if !Authorize(session, route) {
            log.Printf("[%s] Forbidden access by %s to %s -> 403", now, session.Email, r.URL.Path)
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }

        log.Printf("[%s] Access granted for %s (roles: [%s]) to %s", now, session.Email, userRoles, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}


func RequireAuth(route *Route, next http.Handler) http.Handler {
    return Middleware( route, next)
}
