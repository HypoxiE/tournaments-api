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
	DataType string            `json:"data_type"`
	Data     []json.RawMessage `json:"data"`
}

// data_type: "metric", "metadata", "result", "tournament"
// {"data_type": "metric", "data": [{"result_id": 1, "qwerty"}]}
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

	switch data.DataType {
	case "metric":
		for _, jsonMetric := range data.Data {
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
		}
	case "metadata":
		for _, jsonMetadata := range data.Data {
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
	case "result":
		for _, jsonResult := range data.Data {
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
			c.JSON(http.StatusOK, gin.H{"updated": result})
		}
	case "tournament":
		for _, jsonMetadata := range data.Data {
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
	}
}
