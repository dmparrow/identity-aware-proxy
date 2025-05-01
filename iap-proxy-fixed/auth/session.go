
package auth

import (
    "net/http"
    "encoding/base64"
    "encoding/json"
    "strings"
    "log"
)

var sessionCookieName = "iap_session"

type Session struct {
    RawIDToken string
    Roles      []string
    Groups     []string   
    Email      string
    Username   string
}

func GetSession(r *http.Request) (*Session, error) {
    cookie, err := r.Cookie(SessionCookieName)
    if err != nil {
        return nil, err
    }

    roles, email, username := extractRolesAndEmail(cookie.Value)

    return &Session{
        RawIDToken: cookie.Value,
        Roles:      roles,
        Email:      email,
        Username:   username,
    }, nil
}

func (s *Session) IsValid() bool {
    return s.RawIDToken != ""
}

func Authorize(session *Session, route *Route) bool {
    // Check roles
    for _, allowedRole := range route.AllowedRoles {
        for _, role := range session.Roles {
            if role == allowedRole {
                return true
            }
        }
    }

    // Check groups
    for _, allowedGroup := range route.AllowedGroups {
        for _, group := range session.Groups {
            if group == allowedGroup {
                return true
            }
        }
    }

    return false
}

func extractRolesAndEmail(idToken string) ([]string, string, string) {
    parts := strings.Split(idToken, ".")
    if len(parts) != 3 {
        return nil, "", ""
    }

    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return nil, "", ""
    }

    var claims struct {
        RealmAccess struct {
            Roles []string `json:"roles"`
        } `json:"realm_access"`
        Groups            []string `json:"groups"`
        Email             string `json:"email"`
        PreferredUsername string `json:"preferred_username"`
    }

    if err := json.Unmarshal(payload, &claims); err != nil {
        return nil, "", ""
    }

    return claims.RealmAccess.Roles, claims.Email, claims.PreferredUsername
}
