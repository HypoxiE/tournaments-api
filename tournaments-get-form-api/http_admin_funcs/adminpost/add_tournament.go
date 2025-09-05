package adminpost

import (
	"io"
	"log"
	"net/http"
	"tournaments-api/classes"
	"tournaments-api/database"

	"github.com/gin-gonic/gin"
)

func AddTournament(c *gin.Context, manager *database.DataBase) {
	var err error
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournament, err := classes.TournamentDataFromJson(body, c.ClientIP())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := manager.DataBase.Create(&tournament)
	if result.Error != nil {
		log.Fatalf("[ERROR] failed to insert user: %v", result.Error)
	}

	jsonBytes, err := tournament.ToJson()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", jsonBytes)
}
