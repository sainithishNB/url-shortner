package services

import (
	"errors"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sainithishNB/url-shortner.git/cache"
	"github.com/sainithishNB/url-shortner.git/models"
	"github.com/sainithishNB/url-shortner.git/repository"
	"gorm.io/gorm"
)

const base62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var ErrURLExpired = errors.New("url has expired")
var ErrURLNotFound = errors.New("url not found")
var ErrAliasExists = errors.New("alias already exists")

func generateShortCode(db *gorm.DB) string {
	for {
		code := make([]byte, 6)
		for i := range code {
			code[i] = base62[rand.Intn(62)]
		}
		shortCode := string(code)

		var url models.URL
		result := db.Where("short_code=?", shortCode).First(&url)
		if result.Error != nil {
			return shortCode
		}
	}
}
func ShortenURL(db *gorm.DB, rdb *redis.Client, req models.ShortenRequest) (models.URL, error) {
	var shortCode string
	if req.Alias == "" {
		shortCode = generateShortCode(db)
	} else {

		var url models.URL
		result := db.Where("short_code=?", req.Alias).First(&url)
		if result.Error == nil {
			return models.URL{}, ErrAliasExists
		}
		shortCode = req.Alias
	}

	url := models.URL{ShortCode: shortCode, LongURL: req.URL}
	if req.ExpiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)
		url.ExpiresAt = &expiresAt
	}
	err := repository.CreateURL(db, url)
	if err != nil {
		return models.URL{}, err
	}
	ttl := 24 * time.Hour
	if url.ExpiresAt != nil {
		ttl = time.Until(*url.ExpiresAt)
	}
	cache.SetURL(rdb, shortCode, url.LongURL, ttl)
	return url, nil
}
func GetOriginalURL(db *gorm.DB, rdb *redis.Client, code string) (string, error) {
	val, err := cache.GetURL(rdb, code)
	if err == nil {
		url, _ := repository.FindByShortCode(db, code)
		repository.IncrementClickCount(db, url)
		return val, nil
	}
	url, err := repository.FindByShortCode(db, code)
	if err != nil {
		return "", ErrURLNotFound
	}
	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		return "", ErrURLExpired
	}
	ttl := 24 * time.Hour
	if url.ExpiresAt != nil {
		ttl = time.Until(*url.ExpiresAt)
	}
	cache.SetURL(rdb, code, url.LongURL, ttl)
	repository.IncrementClickCount(db, url)

	return url.LongURL, nil

}
func GetStats(db *gorm.DB, code string) (models.URL, error) {
	url, err := repository.FindByShortCode(db, code)
	if err != nil {
		return models.URL{}, ErrURLNotFound
	}
	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		return models.URL{}, ErrURLExpired
	}
	return url, nil
}
