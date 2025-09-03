package main

import (
	"log"
	"tournaments-api/classes"
	"tournaments-api/database"
	"runtime"
	"net/http"
	"io"

	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(1)
	router := gin.Default()

	manager := database.Init()


	router.POST("/registration", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := classes.From_Json(body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		jsonBytes, err := user.To_Json()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result := manager.DataBase.Create(&user)
		if result.Error != nil {
			log.Fatalf("[ERROR] failed to insert user: %v", result.Error)
		}

		log.Printf("[INFO] New registration: %v!", user.Username)

		c.Data(http.StatusOK, "application/json", jsonBytes)
	})

	router.Run(":8080")
}