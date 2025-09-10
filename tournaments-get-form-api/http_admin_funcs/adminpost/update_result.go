package adminpost

import (
	"encoding/json"
	"io"
	"net/http"
	"tournaments-api/database"

	"github.com/gin-gonic/gin"
)

//.Model(&user).Updates(User{Name: "Alice", Age: 30})

type InputData struct {
	DataType string          `json:"data_type"`
	Data     json.RawMessage `json:"data"`
}

// data_type: "metric", "metadata", "result", "tournament"
// {"data_type": "metric", "data": {"result_id": 1, "qwerty"}}
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
	case "metadata":
	case "result":
	case "tournament":
	}
}
