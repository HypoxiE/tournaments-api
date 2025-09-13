package adminget

import (
	"net/http"
	"tournaments-api/classes"
	"tournaments-api/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Leaderbord(c *gin.Context, manager *database.DataBase) {
	tournament_id := c.Query("tournament_id")

	var tournament classes.Tournament
	err := manager.
		DataBase.
		Preload("Results", func(db *gorm.DB) *gorm.DB {
			return db.Order("score DESC, cost DESC")
		}).
		Preload("Results.Metrics", func(db *gorm.DB) *gorm.DB {
			return db.Order("metric_key ASC")
		}).Preload("Results.Metadata", func(db *gorm.DB) *gorm.DB {
		return db.Order("metadata_key ASC")
	}).First(&tournament, tournament_id).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range tournament.Results {
		tournament.Results[i].PublicSteamID = tournament.Results[i].SteamID
		tournament.Results[i].PublicMail = tournament.Results[i].Mail
		tournament.Results[i].PublicIP = tournament.Results[i].IP
	}

	c.JSON(http.StatusOK, tournament)
}
