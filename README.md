# 🛡️ IAP Proxy – Identity-Aware Reverse Proxy with OAuth2

A lightweight Go-based Identity-Aware Proxy that authenticates requests using OAuth2 (OIDC-compliant providers like Keycloak) and reverse proxies to upstream services.

![Go](https://img.shields.io/badge/Go-1.22-blue)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

## 📦 Features

- 🔐 OAuth2 / OpenID Connect (OIDC) login via Keycloak (or any OIDC provider)
- 🔁 Secure reverse proxy to internal services
- 🍪 Session management via signed ID token cookies
- 🧱 Per-route configuration from `config.json`
- 🧑‍💼 Role-based access control
- ⚡ Fast, minimal, self-contained Go binary

---

## 🚀 Getting Started

### 🧱 Requirements

- Go 1.22+
- Docker + Docker Compose

---

### 🐳 Quick Start (Dev)

1. Clone the repo:
   ```bash
   git clone https://github.com/your-user/iap-proxy.git
   cd iap-proxy
   ```

2. Build and run the stack:
   ```bash
   docker compose up --build
   ```

3. Visit the proxy:
   ```http
   http://goapp.localhost
   ```

---

## 🗂 Project Structure

```
├── auth/                   # Core authentication logic
│   ├── config.go           # JSON config loader and route matching
│   ├── handler.go          # /login and /callback handlers
│   ├── middleware.go       # Session validation and role-based auth
│   ├── oauth.go            # OIDC client setup (per-route)
│   └── session.go          # Cookie session parsing and validation
│
├── conf/
│   └── config.json         # Route definitions and OAuth2 settings
│
├── proxy/
│   └── reverse.go          # Reverse proxy handler to upstream services
│
├── Dockerfile              # Container build definition
├── go.mod                  # Go module definition
├── go.sum                  # Dependency hashes
└── main.go      
```

---

## 🔧 `config.json` Example

```json
{
  "routes": [
    {
      "rule": "goapp.localhost",
      "service": "go-app",
      "OAuth2": {
        "provider": {
          "name": "oidc",
          "issuer": "http://keycloak:8080",
          "realm": "dev",
          "clientId": "iap",
          "clientSecret": "your-secret",
          "scopes": ["openid", "profile", "email"],
          "publicAuthHost": "http://goapp.localhost"
        },
        "authResponseHeaders": ["X-Forwarded-User", "X-Forwarded-Email"]
      },
      "Upstream": {
        "passHostHeader": true,
        "servers": [
          { "url": "http://token-echo:8080" }
        ]
      }
    }
  ]
}
```

---

## 🔐 Security Notes

- Secure cookies (`HttpOnly`, `Secure`) should be enabled in production.
- Use TLS on public-facing endpoints.
- Rotate client secrets and restrict allowed redirect URIs in your IdP.

---

## 📊 Metrics & Monitoring

Optional integrations:
- 📈 [Traefik Dashboard](https://grafana.com/grafana/dashboards/12559-traefik-v2-docker/)
- 📉 Node Exporter & cAdvisor for Docker container metrics
- 🔍 Custom app logs via `log.Printf`

---

## 🛠 Future Enhancements

- JWT refresh token support
- Pluggable session storage (Redis, etc.)
- JSON Web Key Set (JWKS) verification
- Multi-client, multi-tenant support
- Advanced routing with path rewrites

---

## 🧑‍💻 Author

**David Arrow**  
[arrowstack.dev](https://arrowstack.dev)

---

## 📜 License

MIT License. Use freely. Contributions welcome!
