package adminpost

import (
	"encoding/json"
	"io"
	"net/http"
	"tournaments-api/classes"
	"tournaments-api/database"

	"github.com/gin-gonic/gin"
)

//.Model(&user).Updates(User{Name: "Alice", Age: 30})

type InputData struct {
	Results     []json.RawMessage `json:"results"`
	Metrics     []json.RawMessage `json:"metrics"`
	Metadata    []json.RawMessage `json:"metadata"`
	Tournaments []json.RawMessage `json:"tournaments"`
}

// data_type: "metric", "metadata", "result", "tournament"
// {"metrics": [], "metadata": [], "results": [], "tournaments": []}
func UpdateResult(c *gin.Context, manager *database.DataBase) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var data InputData
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Metrics
	set := make(map[uint]struct{})
	for _, jsonMetric := range data.Metrics {
		var metric classes.Metric
		if err := json.Unmarshal(jsonMetric, &metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db := manager.DataBase.Model(&metric).Updates(metric)
		if db.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": db.Error})
			return
		}
		if err := db.First(&metric, metric.ID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		set[metric.ResultID] = struct{}{}
	}
	unique := make([]string, 0, len(set))
	for value := range set {
		var result classes.Result
		if err := manager.DataBase.Preload("Tournament").Preload("Metrics").First(&result, value).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result.CalculateScore(result.Tournament)

		db := manager.DataBase.Model(&result).Updates(result)
		if db.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": db.Error})
			return
		}
		unique = append(unique, result.Username)
	}
	c.JSON(http.StatusOK, gin.H{"updated": unique})

	// Metadata
	for _, jsonMetadata := range data.Metadata {
		var metadata classes.Metadata
		if err := json.Unmarshal(jsonMetadata, &metadata); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db := manager.DataBase.Model(&metadata).Updates(metadata)
		if db.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": db.Error})
			return
		}
	}

	// Results
	for _, jsonResult := range data.Results {
		var result classes.Result
		if err := json.Unmarshal(jsonResult, &result); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db := manager.DataBase.Model(&result).Updates(result)
		if db.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": db.Error})
			return
		}

		if err := manager.DataBase.Preload("Tournament").Preload("Metrics").First(&result, result.ID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result.CalculateScore(result.Tournament)

		db = manager.DataBase.Model(&result).Updates(result)
		if db.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": db.Error})
			return
		}

		c.JSON(http.StatusOK, gin.H{"updated": result.Username})
	}

	// Tournaments
	for _, jsonTournament := range data.Tournaments {
		var tournament classes.Tournament
		if err := json.Unmarshal(jsonTournament, &tournament); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db := manager.DataBase.Model(&tournament).Updates(tournament)
		if db.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": db.Error})
			return
		}
	}
}
