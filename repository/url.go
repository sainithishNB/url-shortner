package repository

import (
	"github.com/sainithishNB/url-shortner.git/models"
	"gorm.io/gorm"
)

func FindByShortCode(db *gorm.DB, code string) (models.URL, error) {
	var url models.URL
	result := db.Where("short_code=?", code).First(&url)

	return url, result.Error
}
func CreateURL(db *gorm.DB, url models.URL) error {
	result := db.Create(&url)
	return result.Error
}

func IncrementClickCount(db *gorm.DB, url models.URL) {
	db.Model(&url).Update("click_count", gorm.Expr("click_count + 1"))
}
