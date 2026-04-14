# Copilot Instructions — URL Shortener (Go)

## Project Overview
A URL shortener built in Go using Gin, GORM (MySQL), and Redis. Converts long URLs to short codes and redirects visitors. Refactor from single `main.go` to layered package structure is complete.

## Module Name
```
github.com/sainithishNB/url-shortner.git
```
Use this exact path for all internal imports.

## Architecture (Refactor Complete ✅)
```
main.go              → startup only (router, DB, Redis init)
models/url.go        → URL struct + ShortenRequest struct
config/config.go     → ConnectDB() + ConnectRedis() + AutoMigrate
cache/redis.go       → GetURL(), SetURL(), DelURL()
repository/url.go    → FindByShortCode(), CreateURL(), IncrementClickCount()
services/url.go      → ShortenURL(), GetOriginalURL(), GetStats()
handlers/http.go     → Handler struct, NewHandler(), ShortenHandler(), RedirectHandler(), StatsHandler()
```

## Key Data Model (`models.URL`)
- `ShortCode` → `gorm:"size:10;uniqueIndex"` — must be varchar not text
- `ExpiresAt *time.Time` — pointer = optional, nil means no expiry
- `CreatedAt` — auto-set by GORM, never set manually
- `ClickCount int` — incremented atomically via `gorm.Expr("click_count + 1")`

## Custom Errors (services/url.go)
```go
var ErrURLExpired  = errors.New("url has expired")
var ErrURLNotFound = errors.New("url not found")
```
Handlers check these errors and map to HTTP status codes — services never touch HTTP.

## Cache-Aside Pattern (core flow)
```
GET /:shortCode
  → check Redis first
  → hit: increment MySQL click count, redirect
  → miss: query MySQL, check expiry, set Redis with correct TTL, redirect
```
Redis TTL is always set to `time.Until(*url.ExpiresAt)` if expiry exists, else 24h default.

## URL Validation
URLs must start with `http://` or `https://` — validated in `ShortenHandler` using `strings.HasPrefix`.

## Status Codes
- `201` — short URL created
- `302` — redirect (use StatusFound, never 301)
- `404` — short code not found
- `410` — URL exists but expired
- `400` — invalid input
- `500` — DB save failed

## Local Infrastructure (Podman)
```bash
podman machine start
podman start mysql-url-shortner   # MySQL on port 3307
podman start redis-url-shortner   # Redis on port 6379
```
MySQL DSN: `root:root@tcp(localhost:3307)/urlshortner?parseTime=True`
Redis Addr: `localhost:6379`

## Run & Test
```bash
go run main.go          # starts on :8080
```
Test endpoints with Thunder Client (VS Code) or curl.

## GORM Conventions
- `db.AutoMigrate(&models.URL{})` lives in `config.ConnectDB()` — runs on every startup
- Drop table manually when adding new columns: `DROP TABLE urls;`
- Atomic increment: `db.Model(&url).Update("click_count", gorm.Expr("click_count + 1"))`

## Handler Pattern
```go
// handlers/http.go
type Handler struct {
    db  *gorm.DB
    rdb *redis.Client
}
func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler
// main.go wires: h := handlers.NewHandler(db, rdb)
```
