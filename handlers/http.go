package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sainithishNB/url-shortner.git/models"
	"github.com/sainithishNB/url-shortner.git/services"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{db: db, rdb: rdb}
}

func (h *Handler) ShortenHandler(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": " invalid JSon body"})
		return
	}
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url must start with http:// or https://"})
		return
	}
	url, err := services.ShortenURL(h.db, h.rdb, req)
	if err == services.ErrAliasExists {
		c.JSON(http.StatusConflict, gin.H{"error": "alias already taken"}) // 409!
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"short_code": url.ShortCode,
		"short_url":  "http://localhost:8080/" + url.ShortCode,
		"long_url":   req.URL,
	})
}
func (h *Handler) RedirectHandler(c *gin.Context) {
	code := c.Param("shortCode")
	longURL, err := services.GetOriginalURL(h.db, h.rdb, code)
	if err == services.ErrURLExpired {
		c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
		return
	}
	if err == services.ErrURLNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Redirect(http.StatusFound, longURL)
}
func (h *Handler) StatsHandler(c *gin.Context) {
	code := c.Param("shortCode")
	url, err := services.GetStats(h.db, code)
	if err == services.ErrURLExpired {
		c.JSON(http.StatusGone, gin.H{"error": "URL has expired"})
		return
	}
	if err == services.ErrURLNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"short_code":  code,
		"short_url":   "http://localhost:8080/" + code,
		"long_url":    url.LongURL,
		"created_at":  url.CreatedAt,
		"expires_at":  url.ExpiresAt,
		"click_count": url.ClickCount,
	})
}
