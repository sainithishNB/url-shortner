package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sainithishNB/url-shortner.git/config"
	"github.com/sainithishNB/url-shortner.git/handlers"
	"github.com/sainithishNB/url-shortner.git/middleware"
)

func main() {

	db := config.ConnectDB()
	rdb := config.ConnectRedis()
	h := handlers.NewHandler(db, rdb)
	r := gin.Default()

	r.POST("/shorten", middleware.RateLimit(), h.ShortenHandler)
	r.GET("/:shortCode", h.RedirectHandler)
	r.GET("/:shortCode/stats", h.StatsHandler)
	r.Run(":8080")
}
