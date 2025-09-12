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

type ResponseData struct {
	Results       []string
	ResultsID     []uint
	MetricsID     []uint
	MetadataID    []uint
	Tournaments   []string
	TournamentsID []uint
}

func (response *ResponseData) unique() {
	response.Results = classes.GetUnique(response.Results)
	response.ResultsID = classes.GetUnique(response.ResultsID)
	response.MetadataID = classes.GetUnique(response.MetadataID)
	response.MetricsID = classes.GetUnique(response.MetricsID)
	response.Tournaments = classes.GetUnique(response.Tournaments)
	response.TournamentsID = classes.GetUnique(response.TournamentsID)
}

type UnpackMetricsOutput struct {
	MetricsID   []uint
	ResultsID   []uint
	ResultsName []string
}

func UnpackMetrics(manager *database.DataBase, metrics []json.RawMessage) (UnpackMetricsOutput, error) {
	var output UnpackMetricsOutput

	set := make(map[uint]struct{})
	for _, jsonMetric := range metrics {
		var metric classes.Metric
		if err := json.Unmarshal(jsonMetric, &metric); err != nil {
			return output, err
		}
		db := manager.DataBase.Model(&metric).Updates(metric)
		if db.Error != nil {
			return output, db.Error
		}
		if err := db.First(&metric, metric.ID).Error; err != nil {
			return output, db.Error
		}
		set[metric.ResultID] = struct{}{}
		output.MetricsID = append(output.MetricsID, metric.ID)
	}
	for value := range set {
		var result classes.Result
		if err := manager.DataBase.Preload("Tournament").Preload("Metrics").First(&result, value).Error; err != nil {
			return output, err
		}
		result.CalculateScore(result.Tournament)

		db := manager.DataBase.Model(&result).Updates(result)
		if db.Error != nil {
			return output, db.Error
		}

		output.ResultsName = append(output.ResultsName, result.Username)
		output.ResultsID = append(output.ResultsID, result.ID)
	}
	return output, nil
}

type UnpackMetadataOutput struct {
	MetadataID []uint
}

func UnpackMetadata(manager *database.DataBase, metadata []json.RawMessage) (UnpackMetadataOutput, error) {
	var output UnpackMetadataOutput
	for _, jsonMetadata := range metadata {
		var metadata classes.Metadata
		if err := json.Unmarshal(jsonMetadata, &metadata); err != nil {
			return output, err
		}
		db := manager.DataBase.Model(&metadata).Updates(metadata)
		if db.Error != nil {
			return output, db.Error
		}
		output.MetadataID = append(output.MetadataID, metadata.ID)
	}
	return output, nil
}

type UnpackResultsOutput struct {
	ResultsID   []uint
	ResultsName []string
}

func UnpackResults(manager *database.DataBase, results []json.RawMessage) (UnpackResultsOutput, error) {
	var output UnpackResultsOutput

	for _, jsonResult := range results {
		var result classes.Result
		if err := json.Unmarshal(jsonResult, &result); err != nil {
			return output, err
		}
		db := manager.DataBase.Model(&result).Updates(result)
		if db.Error != nil {
			return output, db.Error
		}

		if err := manager.DataBase.Preload("Tournament").Preload("Metrics").First(&result, result.ID).Error; err != nil {
			return output, err
		}
		result.CalculateScore(result.Tournament)
		result.PublicToPrivate()

		db = manager.DataBase.Model(&result).Updates(result)
		if db.Error != nil {
			return output, db.Error
		}
		output.ResultsName = append(output.ResultsName, result.Username)
		output.ResultsID = append(output.ResultsID, result.ID)
	}
	return output, nil
}

type UnpackTournamentsOutput struct {
	TournamentsID   []uint
	TournamentsName []string
	ResultsID       []uint
	ResultsName     []string
}

func UnpackTournaments(manager *database.DataBase, tournaments []json.RawMessage) (UnpackTournamentsOutput, error) {
	var output UnpackTournamentsOutput

	for _, jsonTournament := range tournaments {
		var tournament classes.Tournament
		if err := json.Unmarshal(jsonTournament, &tournament); err != nil {
			return output, err
		}
		db := manager.DataBase.Model(&tournament).Updates(tournament)
		if db.Error != nil {
			return output, db.Error
		}
		output.TournamentsID = append(output.TournamentsID, tournament.ID)
		output.TournamentsName = append(output.TournamentsName, tournament.Name)
	}

	return output, nil
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

	var response ResponseData

	// Metrics
	if output, err := UnpackMetrics(manager, data.Metrics); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {
		response.MetricsID = append(response.MetricsID, output.MetricsID...)
		response.ResultsID = append(response.ResultsID, output.ResultsID...)
		response.Results = append(response.Results, output.ResultsName...)
	}

	// Metadata
	if output, err := UnpackMetadata(manager, data.Metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {
		response.MetadataID = append(response.MetadataID, output.MetadataID...)
	}

	// Results
	if output, err := UnpackResults(manager, data.Results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {
		response.ResultsID = append(response.ResultsID, output.ResultsID...)
		response.Results = append(response.Results, output.ResultsName...)
	}

	// Tournaments
	if output, err := UnpackTournaments(manager, data.Tournaments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {
		response.Tournaments = append(response.Tournaments, output.TournamentsName...)
		response.TournamentsID = append(response.TournamentsID, output.TournamentsID...)
		response.Results = append(response.Results, output.ResultsName...)
		response.ResultsID = append(response.ResultsID, output.ResultsID...)
	}

	response.unique()
	c.JSON(http.StatusOK, gin.H{"updated": response})
}
