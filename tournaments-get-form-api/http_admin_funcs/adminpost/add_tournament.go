package adminpost

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"tournaments-api/classes"
	"tournaments-api/database"

	"github.com/gin-gonic/gin"
)

// {"name": "Test Tournament", "stop_timestamp": 1757099263, "metadata": ["streams", "comment"], "variables": ["animals", "humans"], "formula": "(animals * 2) + (humans * 3)"}
func TournamentDataFromJson(jsonData []byte, ip string) (classes.Tournament, error) {
	var tournament classes.Tournament
	if err := json.Unmarshal(jsonData, &tournament); err != nil {
		return tournament, err
	}

	if tournament.Name == "" {
		err := errors.New("error: The name field is missing")
		return tournament, err
	} else if tournament.StopTimestamp == 0 {
		err := errors.New("error: The stop_timestamp field is missing")
		return tournament, err
	}

	tournament.ID = 0

	return tournament, nil
}

func AddTournament(c *gin.Context, manager *database.DataBase) {
	var err error
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tournament, err := TournamentDataFromJson(body, c.ClientIP())
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
