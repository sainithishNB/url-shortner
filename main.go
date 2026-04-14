package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sainithishNB/url-shortner.git/config"
	"github.com/sainithishNB/url-shortner.git/handlers"
)

func main() {

	db := config.ConnectDB()
	rdb := config.ConnectRedis()
	h := handlers.NewHandler(db, rdb)
	r := gin.Default()

	r.POST("/shorten", h.ShortenHandler)
	r.GET("/:shortCode", h.RedirectHandler)
	r.GET("/:shortCode/stats", h.StatsHandler)
	r.Run(":8080")
}
