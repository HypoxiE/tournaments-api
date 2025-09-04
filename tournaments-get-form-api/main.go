package main

import (
	"log"
	"tournaments-api/database"
	"tournaments-api/http_funcs/post"
	"runtime"

	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(1)
	router := gin.Default()

	manager := database.Init()
	log.Println("[INFO] Инициализация завершена!")

	router.POST("/registration", func(c *gin.Context) {
		post.Registration(c, &manager)
	})

	router.Run(":8080")
}