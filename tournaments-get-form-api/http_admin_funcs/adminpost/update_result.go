package adminpost

import (
	"tournaments-api/database"

	"github.com/gin-gonic/gin"
)

// data_type: "metric", "metadata", "result"
// {"result_id": 4, "tournament_id": 3, "data_type": "result", "update_data": [{"key": "version"}]}
func UpdateResult(c *gin.Context, manager *database.DataBase) {

}
