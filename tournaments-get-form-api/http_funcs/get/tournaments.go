package get

import (
	"net/http"
	"tournaments-api/classes"
	"tournaments-api/database"

	"github.com/gin-gonic/gin"
)

type OutputTournaments struct {
	classes.Tournament
	NumberResults uint `json:"number_results"`
}

func Tournaments(c *gin.Context, manager *database.DataBase) {

	var tournaments []OutputTournaments

	if err := manager.DataBase.Model(&classes.Tournament{}).
		Select("tournaments.*, COUNT(results.result_id) AS number_results").
		Joins("LEFT JOIN results ON results.tournament_id = tournaments.tournament_id").
		Group("tournaments.tournament_id").Scan(&tournaments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tournaments)
}
