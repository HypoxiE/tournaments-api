package main

import (
	"log"
	"tournaments-api/database"
	"tournaments-api/http_funcs/post"
	"tournaments-api/http_funcs/get"
	"net/http"
	"runtime"

	"golang.org/x/time/rate"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b) // r = requests per second, b = burst

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	router := gin.Default()
	router.Use(RateLimitMiddleware(5, 10))

	manager := database.Init()
	log.Println("[INFO] Инициализация завершена!")

	router.POST("/registration", func(c *gin.Context) {
		post.Registration(c, &manager)
	})

	router.GET("/leaderbord", func(c *gin.Context) {
		get.Leaderbord(c, &manager)
	})

	router.Run(":8080")
}