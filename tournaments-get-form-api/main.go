package main

import (
	"log"
	"net/http"
	"runtime"
	"tournaments-api/database"
	"tournaments-api/http_admin_funcs/adminget"
	"tournaments-api/http_admin_funcs/adminpost"
	"tournaments-api/http_funcs/get"
	"tournaments-api/http_funcs/post"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
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
	router.Use(cors.Default())
	router.Use(RateLimitMiddleware(5, 10))

	adminRouter := gin.Default()
	adminRouter.Use(cors.Default())

	manager := database.Init()
	log.Println("[INFO] Инициализация завершена!")

	router.POST("/registration", func(c *gin.Context) {
		post.Registration(c, &manager)
	})

	router.GET("/leaderboard", func(c *gin.Context) {
		get.Leaderbord(c, &manager)
	})

	adminRouter.POST("/admin/add_tournament", func(c *gin.Context) {
		adminpost.AddTournament(c, &manager)
	})

	adminRouter.POST("/admin/leaderboard", func(c *gin.Context) {
		adminget.Leaderbord(c, &manager)
	})

	go func() {
		if err := adminRouter.Run(":9090"); err != nil {
			panic(err)
		}
	}()

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
