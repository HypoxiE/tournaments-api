package get

import (
	"net/http"
	"strconv"
	"tournaments-api/classes"
	"tournaments-api/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Leaderbord(c *gin.Context, manager *database.DataBase) {
	tournament_id_str := c.Query("tournament_id")
	tournament_id, err := strconv.Atoi(tournament_id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result classes.Tournament
	if err := manager.
		DataBase.
		Preload("Results", func(db *gorm.DB) *gorm.DB {
			return db.Where("status >= ?", 0).Order("score DESC, cost DESC")
		}).
		Preload("Results.Metrics").Preload("Results.Metadata").First(&result, tournament_id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
