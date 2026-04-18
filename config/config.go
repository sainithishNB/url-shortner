package config

import (
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/sainithishNB/url-shortner.git/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := os.Getenv("MYSQL_URL")
	if dsn == "" {
		dsn = "root:root@tcp(localhost:3307)/urlshortner?parseTime=True"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&models.URL{})
	return db
}
func ConnectRedis() *redis.Client {
	addr := os.Getenv("REDIS_URL")
	if addr == "" {
        // local development fallback
        addr = "localhost:6379"
    }
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
