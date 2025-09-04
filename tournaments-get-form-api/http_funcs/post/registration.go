package post

import (
	"log"
	"fmt"
	"tournaments-api/classes"
	"tournaments-api/database"
	"net/http"
	"io"

	"github.com/gin-gonic/gin"
)

func Registration(c *gin.Context, manager *database.DataBase) {
	var err error
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := classes.From_Json(body, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tournament classes.Tournament
	err = manager.DataBase.First(&tournament, user.TournamentID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	find_key := func(metric classes.Metric) bool {
		for _, variable := range tournament.Variables {
			if variable == metric.Key {
				return true
			}
		}
		return false
	}
	for _, metric := range user.Metrics {
		if !find_key(metric) {
			err = fmt.Errorf("Error: unknown metric %v", metric.Key)
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = user.CalculateScore(tournament)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
}